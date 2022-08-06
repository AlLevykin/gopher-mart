-- +goose Up
CREATE SCHEMA IF NOT EXISTS gophermart;
CREATE TABLE "user" (
                      login varchar(20) PRIMARY KEY,
                      pwh text,
                      salt text
);
-- +goose Down
DROP TABLE "user";