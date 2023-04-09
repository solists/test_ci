-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists tokens (
   id bigserial PRIMARY KEY,
   token TEXT NOT NULL,
   source TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tokens;
-- +goose StatementEnd
