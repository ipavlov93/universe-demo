-- +goose Up
-- +goose StatementBegin
create table IF NOT EXISTS products
(
    id           serial primary key,
    name varchar(400) unique     not null,
    description  text null,
    created_at   timestamp default now() not null
);
-- +goose StatementEnd