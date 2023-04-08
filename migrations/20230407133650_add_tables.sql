-- +goose Up
-- +goose StatementBegin
create table if not exists ids (
    id bigserial primary key,
    created_at timestamptz not null default now(),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ids;
-- +goose StatementEnd
