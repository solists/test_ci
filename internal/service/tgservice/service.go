package tgservice

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/solists/test_ci/pkg/logger"
	"io"
	openai2 "mymod/internal/client/openai"
	"mymod/internal/controller"
	"mymod/internal/models/openai"
	repomodels "mymod/internal/models/repository"
	"mymod/internal/repository"
	"net/http"
	"os"
	"strings"
)

const defaultLogLimit = 10

var (
	queryRequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_qet_query_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"user"},
	)
)

type Service struct {
	ctrl       controller.IController
	repo       repository.IRepository
	downloader *Downloader
}

func NewService(
	repo repository.IRepository,
	ctrl controller.IController,
	downloader *Downloader,
) *Service {
	return &Service{
		repo:       repo,
		ctrl:       ctrl,
		downloader: downloader,
	}
}

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
		inlineQuery := update.InlineQuery.Query

		queryEndSuffix := "!!!"
		if !strings.HasSuffix(inlineQuery, queryEndSuffix) {
			return
		} else {
			inlineQuery = strings.TrimSuffix(inlineQuery, queryEndSuffix)
		}

		user, err := s.repo.GetUserData(ctx, updateUserFrom.ID)
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

		if user.UserID == 0 || !user.Allowed {
			logger.Infof("user not allowed: %v", user.UserID)
			return
		}

		openaiResp, err := s.ctrl.GetQuery(ctx, &openai.GetQueryRequest{
			UserID: user.UserID,
			Messages: []openai.PromptMessage{
				{
					inlineQuery,
				},
			},
		})
		_, err = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results: []models.InlineQueryResult{
				&models.InlineQueryResultArticle{
					ID:    "1",
					Title: "gpt resp",
					InputMessageContent: models.InputTextMessageContent{
						MessageText:           openaiResp.Result,
						DisableWebPagePreview: true,
					},
				},
			},
		})

		if err != nil {
			logger.Errorf("ohh epic fail inline inlineQuery:  %v, user: %v", err, user)
		}

		return
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
		logger.Info(gotVoice.FilePath)
		fileReader, err := s.downloader.Download(ctx, gotVoice.FilePath)
		if err != nil {
			logger.Errorf("err download voice: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}

		file, err := os.CreateTemp("", "voice-*.mp3")
		if err != nil {
			logger.Errorf("err create temp: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}
		defer func() {
			os.Remove(file.Name())
		}()

		_, err = io.Copy(file, fileReader)
		if err != nil {
			file.Close()
			logger.Errorf("err copy to file voice: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}
		file.Close()

		transcript, err := s.ctrl.GetTranscription(ctx, &openai.GetTranscriptionRequest{
			UserID:   user.UserID,
			FilePath: file.Name(),
		})
		if err != nil {
			logger.Errorf("err GetTranscription: %v, user: %v", err, user)
			finalErr = errors.New("error occurred")
			return
		}

		updateMessage = transcript.Result
	}

	messageLogsIns = append(messageLogsIns, repomodels.MessageLog{
		UserID:    &updateUserFrom.ID,
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Message:   updateMessage,
	})

	var messages []openai.PromptMessage
	// we get from query in reverse order, so first is the last
	for i := len(messageLogs) - 1; i >= 0; i-- {
		messages = append(messages, openai.PromptMessage{
			Message: messageLogs[i].Message,
		})
	}
	messages = append(messages, openai.PromptMessage{
		Message: updateMessage,
	})

	openaiResp, err := s.ctrl.GetQuery(ctx, &openai.GetQueryRequest{
		UserID:   user.UserID,
		Messages: messages,
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

type Downloader struct {
	token string

	server      string
	downloadURL string
	httpClient  Doer
}

func NewDownloader(token string) *Downloader {
	return &Downloader{
		token:       token,
		server:      "https://api.telegram.org",
		downloadURL: "%s/file/bot%s/%s",
		httpClient:  http.DefaultClient,
	}
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

func (d *Downloader) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	url := d.buildDownloadURL(d.token, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	res, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		return nil, fmt.Errorf("while download: %v, %v", res.StatusCode, res.Status)
	}

	return res.Body, nil
}

func (d *Downloader) buildDownloadURL(token, path string) string {
	return fmt.Sprintf(d.downloadURL, d.server, token, path)
}
