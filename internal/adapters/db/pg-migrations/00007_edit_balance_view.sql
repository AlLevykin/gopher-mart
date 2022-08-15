-- +goose Up
DROP VIEW balance;
CREATE VIEW balance AS
SELECT "user", SUM(accrual) - SUM(COALESCE("sum",0)) AS "current", SUM(COALESCE("sum",0)) AS withdrawn FROM "order" LEFT OUTER JOIN withdraw ON "number" = "order" GROUP BY "user";
-- +goose Down
DROP VIEW balance;
CREATE VIEW balance AS
SELECT "user", SUM(accrual) AS "current", SUM(COALESCE("sum",0)) AS withdrawn FROM "order" LEFT OUTER JOIN withdraw ON "number" = "order" GROUP BY "user";