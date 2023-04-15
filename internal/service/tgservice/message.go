package tgservice

import (
	"context"
	"database/sql"
	"fmt"
	openai2 "mymod/internal/client/openai"
	"mymod/internal/models/openai"
	repomodels "mymod/internal/models/repository"
	"os"
	"unicode/utf8"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/solists/test_ci/pkg/logger"
)

func (s *Service) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var finalErr error
	defer func() {
		if finalErr != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   finalErr.Error(),
			})
			if err != nil {
				logger.Errorf("err send message: chatId: %v, error: %v", update.Message.Chat.ID, err)
			}
		}
	}()

	var isInline bool
	if update == nil || update.Message == nil || update.Message.From == nil {
		if update.InlineQuery == nil || update.InlineQuery.From == nil {
			return
		}
		isInline = true
	}

	var updateUserFrom *models.User
	if isInline {
		updateUserFrom = update.InlineQuery.From
	} else {
		updateUserFrom = update.Message.From
	}

	queryRequestCounter.With(prometheus.Labels{"user": fmt.Sprint(updateUserFrom.ID)}).Inc()

	if isInline {
		if err := s.handleInline(ctx, updateUserFrom,
			update.InlineQuery.Query, b, update.InlineQuery.ID); err != nil {
			logger.Errorf("handleInline: %v", err)
		}
	}

	var user repomodels.UserData
	var messageLogsIns []repomodels.MessageLog
	defer func() {
		if err := s.repo.InsertMessageLogs(ctx, messageLogsIns); err != nil {
			logger.Errorf("InsertMessageLogs: %v, messages: %v", err, messageLogsIns)
		}
	}()

	messageLogs, err := s.repo.GetMessageLogWithUserData(ctx, update.Message.Chat.ID, defaultLogLimit)
	if err != nil {
		logger.Errorf("error GetMessageLogWithUserData: %v, chatID: %v", err, update.Message.Chat.ID)
		finalErr = errors.New("error occurred")
		return
	}
	if len(messageLogs) == 0 {
		userResp, err := s.repo.GetUserData(ctx, updateUserFrom.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				userIns := &repomodels.UserData{
					UserID:    updateUserFrom.ID,
					ChatID:    update.Message.Chat.ID,
					FirstName: updateUserFrom.FirstName,
					LastName:  updateUserFrom.LastName,
					UserName:  updateUserFrom.Username,
				}
				if err = s.repo.InsertUserData(ctx, userIns); err != nil {
					logger.Errorf("error InsertUserData: %v, user: %v", err, user)
				}

				finalErr = errors.New("You are not registered yet")
				return
			}
			logger.Errorf("error GetUserData: %v, userID: %v", err, updateUserFrom.ID)

			finalErr = errors.New("error occurred")
			return
		}
		user = *userResp
	}
	if len(messageLogs) > 0 {
		if messageLogs[0].UserID != nil {
			user.UserID = *messageLogs[0].UserID
			user.Allowed = *messageLogs[0].Allowed
			user.ChatID = messageLogs[0].ChatID
		} else {
			userIns := &repomodels.UserData{
				UserID:    updateUserFrom.ID,
				ChatID:    update.Message.Chat.ID,
				FirstName: updateUserFrom.FirstName,
				LastName:  updateUserFrom.LastName,
				UserName:  updateUserFrom.Username,
			}
			if err = s.repo.InsertUserData(ctx, userIns); err != nil {
				logger.Errorf("error InsertUserData: %v, user: %v", err, user)
			}

			finalErr = errors.New("You are not registered yet")
			return
		}
	}

	if user.UserID == 0 || !user.Allowed {
		logger.Infof("user not allowed: %v", user.UserID)
		finalErr = errors.New("You are not registered yet")
		return
	}
	if user.ChatID != update.Message.Chat.ID {
		if err = s.repo.UpdateUserDataChatID(ctx, update.Message.Chat.ID, user.UserID); err != nil {
			logger.Errorf("error UpdateUserDataChatID: %v, user: %v, chat: %v",
				err, user, update.Message.Chat.ID)

			finalErr = errors.New("error occurred")
			return
		}
		user.ChatID = update.Message.Chat.ID
	}

	updateMessage := update.Message.Text
	if update.Message.Voice != nil {
		gotVoice, err := b.GetFile(ctx, &bot.GetFileParams{FileID: update.Message.Voice.FileID})
		if err != nil {
			logger.Errorf("err get voice: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}

		voiceFilenameMP3, err := s.downloader.DownloadVoiceMP3(ctx, gotVoice.FilePath)
		defer func() {
			if voiceFilenameMP3 != nil {
				_ = os.Remove(*voiceFilenameMP3)
			}
		}()
		if err != nil {
			logger.Errorf("err DownloadVoiceMP3: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}
		transcript, err := s.ctrl.GetTranscription(ctx, &openai.GetTranscriptionRequest{
			UserID:   user.UserID,
			FilePath: *voiceFilenameMP3,
		})
		if err != nil {
			logger.Errorf("err GetTranscription: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}
		logger.Infof(transcript.Result)

		updateMessage = transcript.Result
	}

	messageLogsIns = append(messageLogsIns, repomodels.MessageLog{
		UserID:    &updateUserFrom.ID,
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Message:   updateMessage,
	})

	openaiResp, err := s.ctrl.GetQuery(ctx, &openai.GetQueryRequest{
		UserID:   user.UserID,
		Messages: s.createRequestMessages(messageLogs, updateMessage),
	})
	if err != nil {
		logger.Errorf("GetQuery: %v, user: %v", err, user)

		if errors.Is(err, openai2.ErrTooBigPrompt) {
			finalErr = err
		}
		return
	}

	messageSent, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   openaiResp.Result,
	})

	var userIDMessageResp *int64
	if messageSent.From != nil {
		userIDMessageResp = &messageSent.From.ID
	}
	messageLogsIns = append(messageLogsIns, repomodels.MessageLog{
		UserID:    userIDMessageResp,
		ChatID:    messageSent.Chat.ID,
		MessageID: messageSent.ID,
		Message:   messageSent.Text,
	})

	if err != nil {
		logger.Errorf("ohh epic fail:  %v, user: %v", err, user)
	}
}

// euristic size about 2.5 runes per token
const maxQuerySize = 5000

func (s *Service) createRequestMessages(messageLogs []repomodels.MessageLogUserIDResp, updateMessage string) (messages []openai.PromptMessage) {
	var currentRuneCount = utf8.RuneCountInString(updateMessage)
	// we get from query in reverse order, so first is the last
	for i := len(messageLogs) - 1; i >= 0; i-- {
		currentRuneCount += utf8.RuneCountInString(messageLogs[i].Message)
		if currentRuneCount > maxQuerySize {
			break
		}
		messages = append(messages, openai.PromptMessage{
			Message: messageLogs[i].Message,
		})
	}
	messages = append(messages, openai.PromptMessage{
		Message: updateMessage,
	})

	return
}
