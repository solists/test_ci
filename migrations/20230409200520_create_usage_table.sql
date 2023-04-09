-- +goose Up
-- +goose StatementBegin
create table if not exists usage
(
    id         bigserial primary key,
    user_id    bigint unique not null ,
    used_prompt bigint not null default 0,
    used_completed bigint not null default 0,
    used_total bigint not null default 0
);

create unique index usage_user_id_idx on usage(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists usage;
-- +goose StatementEnd
