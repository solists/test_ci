package tgservice

import (
	"context"
	"database/sql"
	"fmt"
	"mymod/internal/models/openai"
	repomodels "mymod/internal/models/repository"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/solists/test_ci/pkg/logger"
)

func (s *Service) handleInline(ctx context.Context,
	updateUserFrom *models.User,
	inlineQuery string,
	b *bot.Bot,
	queryID string,
) error {
	queryEndSuffix := "!!!"
	if !strings.HasSuffix(inlineQuery, queryEndSuffix) {
		return nil
	} else {
		inlineQuery = strings.TrimSuffix(inlineQuery, queryEndSuffix)
	}

	user, err := s.repo.GetUserData(ctx, updateUserFrom.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userIns := &repomodels.UserData{
				UserID:    updateUserFrom.ID,
				ChatID:    192,
				FirstName: updateUserFrom.FirstName,
				LastName:  updateUserFrom.LastName,
				UserName:  updateUserFrom.Username,
			}
			if err = s.repo.InsertUserData(ctx, userIns); err != nil {
				return fmt.Errorf("error InsertUserData: %v, user: %v", err, user)
			}

			return nil
		}

		return fmt.Errorf("error GetUserData: %v, userID: %v", err, updateUserFrom.ID)
	}

	if user.UserID == 0 || !user.Allowed {
		logger.Infof("user not allowed: %v", user.UserID)
		return nil
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
		InlineQueryID: queryID,
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
		return fmt.Errorf("ohh epic fail inline inlineQuery:  %v, user: %v", err, user)
	}

	return nil
}
