package tgservice

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	openai2 "mymod/internal/client/openai"
	"mymod/internal/controller"
	"mymod/internal/models/openai"
	repomodels "mymod/internal/models/repository"
	"mymod/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/solists/test_ci/pkg/logger"
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
	ctrl controller.IController
	repo repository.IRepository
}

func NewService(
	repo repository.IRepository,
	ctrl controller.IController,
) *Service {
	return &Service{
		repo: repo,
		ctrl: ctrl,
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

	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	queryRequestCounter.With(prometheus.Labels{"user": fmt.Sprint(update.Message.From.ID)}).Inc()

	var user repomodels.UserData
	var messageLogsIns []repomodels.MessageLog
	defer func() {
		if err := s.repo.InsertMessageLogs(ctx, messageLogsIns); err != nil {
			logger.Errorf("InsertMessageLogs: %v, messages: %v", err, messageLogsIns)
		}
	}()

	var userIDMessageReq *int64
	if update.Message.From != nil {
		userIDMessageReq = &update.Message.From.ID
	}
	messageLogsIns = append(messageLogsIns, repomodels.MessageLog{
		UserID:    userIDMessageReq,
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Message:   update.Message.Text,
	})

	messageLogs, err := s.repo.GetMessageLogWithUserData(ctx, update.Message.Chat.ID, defaultLogLimit)
	if err != nil {
		logger.Errorf("error GetMessageLogWithUserData: %v, chatID: %v", err, update.Message.Chat.ID)
		finalErr = errors.New("error occurred")
		return
	}
	if len(messageLogs) == 0 {
		userResp, err := s.repo.GetUserData(ctx, update.Message.From.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				userIns := &repomodels.UserData{
					UserID:    update.Message.From.ID,
					ChatID:    update.Message.Chat.ID,
					FirstName: update.Message.From.FirstName,
					LastName:  update.Message.From.LastName,
					UserName:  update.Message.From.Username,
				}
				if err = s.repo.InsertUserData(ctx, userIns); err != nil {
					logger.Errorf("error InsertUserData: %v, user: %v", err, user)
				}

				finalErr = errors.New("You are not registered yet")
				return
			}
			logger.Errorf("error GetUserData: %v, userID: %v", err, update.Message.From.ID)

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
				UserID:    update.Message.From.ID,
				ChatID:    update.Message.Chat.ID,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
				UserName:  update.Message.From.Username,
			}
			if err = s.repo.InsertUserData(ctx, userIns); err != nil {
				logger.Errorf("error InsertUserData: %v, user: %v", err, user)
			}

			finalErr = errors.New("You are not registered yet")
			return
		}
	}
	if user.UserID == 0 || !user.Allowed {
		logger.Infof("user not allowed: %v", update.Message.From.ID)
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

	var messages []openai.PromptMessage
	// we get from query in reverse order, so first is the last
	for i := len(messageLogs) - 1; i >= 0; i-- {
		messages = append(messages, openai.PromptMessage{
			Message: messageLogs[i].Message,
		})
	}
	messages = append(messages, openai.PromptMessage{
		Message: update.Message.Text,
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
