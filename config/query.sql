-- name: ListCompanies :many
SELECT *
FROM company;

-- name: SaveZacksDailyRow :exec
INSERT INTO zacks_daily (
        symbol,
        company,
        price,
        dollar_change,
        percent_change,
        industry_rank,
        zacks_rank,
        value_score,
        growth_score,
        momentum_score,
        vgm_score
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
