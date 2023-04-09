package repository

type Usage struct {
	ID            uint64 `db:"id" json:"id"`
	UserID        uint64 `db:"user_id" json:"user_id" validate:"required"`
	UsedPrompt    uint64 `db:"used_prompt" json:"used_prompt"`
	UsedCompleted uint64 `db:"used_completed" json:"used_completed"`
	UsedTotal     uint64 `db:"used_total" json:"used_total"`
}

type UsageInsert struct {
	UserID        uint64 `db:"user_id" json:"user_id" validate:"required"`
	UsedPrompt    uint64 `db:"used_prompt" json:"used_prompt"`
	UsedCompleted uint64 `db:"used_completed" json:"used_completed"`
	UsedTotal     uint64 `db:"used_total" json:"used_total"`
}
