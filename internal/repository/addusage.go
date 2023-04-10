package repository

import (
	"context"
	"mymod/internal/models/repository"
)

func (r *Repository) AddUsage(ctx context.Context, req *repository.UsageInsert) (err error) {
	if _, err = r.db.NamedExecContext(ctx, query, req); err != nil {
		return
	}

	return
}

const query = `
INSERT INTO usage (user_id, used_prompt, used_completed, used_total)
VALUES (:user_id, :used_prompt, :used_completed, :used_total)
ON CONFLICT (user_id) DO UPDATE SET
    used_prompt = excluded.used_prompt + usage.used_prompt,
    used_completed = excluded.used_completed + usage.used_completed,
    used_total = excluded.used_total + usage.used_total
`
