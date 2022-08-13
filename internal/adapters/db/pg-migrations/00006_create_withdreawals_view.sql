-- +goose Up
CREATE VIEW withdrawals AS
SELECT "order", "sum", "processed" AS "processed_at", "order"."user" FROM withdraw LEFT JOIN "order" ON withdraw."order" = "order".number;
-- +goose Down
DROP VIEW withdrawals;