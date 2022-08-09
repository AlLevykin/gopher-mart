-- +goose Up
CREATE TYPE orderType AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE "order" (
                    "number" varchar(12) PRIMARY KEY,
                    "user" varchar(20),
                    status orderType,
                    uploaded date
);
CREATE UNIQUE INDEX order_user_number
    ON "order" ("user", "number");
-- +goose Down
DROP INDEX order_user_number;
DROP TABLE "order";
DROP TYPE orderType;