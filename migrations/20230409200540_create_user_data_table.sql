-- +goose Up
-- +goose StatementBegin
create table if not exists user_data
(
    id          bigserial primary key,
    user_id     bigint unique not null,
    allowed     boolean not null default false,
    chat_id     bigint not null,
    first_name  text,
    last_name   text,
    user_name   text not null default '',
    created_at  timestamptz not null default now()
);

create unique index user_data_chat_id_idx on user_data(user_id, chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_data;
-- +goose StatementEnd
