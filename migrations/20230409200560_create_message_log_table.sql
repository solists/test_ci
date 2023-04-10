-- +goose Up
-- +goose StatementBegin
create table if not exists message_log
(
    id          bigserial primary key,
    user_id     bigint,
    chat_id     bigint not null,
    message_id  bigint not null,
    message     text not null,
    created_at  timestamptz not null default now()
);

create unique index message_data_chat_id_idx on message_log(chat_id, message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists message_log;
-- +goose StatementEnd
