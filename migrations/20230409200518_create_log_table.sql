-- +goose Up
-- +goose StatementBegin
create table if not exists audit_log
(
    id         bigserial primary key,
    user_id    bigint,
    data       jsonb       not null,
    created_at timestamptz not null default now(),
    operation  text        not null,
    response   jsonb,
    status     integer
);

create index user_id_idx on audit_log(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists audit_log;
-- +goose StatementEnd
