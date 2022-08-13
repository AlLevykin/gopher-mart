-- +goose Up
CREATE VIEW balance AS
SELECT "user", SUM(accrual) AS "current", SUM(COALESCE("sum",0)) AS withdrawn FROM "order" LEFT OUTER JOIN withdraw ON "number" = "order" GROUP BY "user";
-- +goose Down
DROP VIEW balance;