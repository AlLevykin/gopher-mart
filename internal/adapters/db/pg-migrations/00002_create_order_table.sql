-- +goose Up
CREATE TYPE orderType AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE "order" (
                    "number" varchar(20) PRIMARY KEY,
                    "user" varchar(20),
                    status orderType,
                    accrual numeric DEFAULT 0,
                    uploaded timestamp DEFAULT NOW()
);
CREATE UNIQUE INDEX order_user_number
    ON "order" ("user", "number");
-- +goose Down
DROP INDEX order_user_number;
DROP TABLE "order";
DROP TYPE orderType;