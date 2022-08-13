-- +goose Up
CREATE TABLE withdraw (
                    "order" varchar(20),
                    "sum" numeric
);
-- +goose Down
DROP TABLE withdraw;