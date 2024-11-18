-- name: ListCompanies :many
SELECT *
FROM company;
-- name: SaveTaskHistory :exec
INSERT INTO task_history (
        task_name,
        task_status,
        start_time,
        end_time,
        details
    )
VALUES (?, ?, ?, ?, ?);
-- name: GetLatestTaskHistory :one
SELECT *
FROM task_history
WHERE task_name = ?
ORDER BY start_time DESC
LIMIT 1;
-- name: SaveZacksDailyRow :exec
INSERT INTO zacks_daily (
        created,
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: SaveZacksGrowthRow :exec
INSERT INTO zacks_growth (
        created,
        symbol,
        company,
        price,
        growth_score,
        year_over_year_q0_growth,
        long_term_growth_percent,
        last_financial_year_actual,
        next_finanical_year_est,
        this_financial_year_est,
        q1_est,
        earnings_expected_surprise_prediction,
        next_report_date
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: SaveYahooInsightsRow :exec
INSERT INTO yahoo_insights (
        created,
        symbol,
        company_name,
        short_term,
        mid_term,
        long_term,
        estimated_return,
        fair_value
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?);
-- name: SaveYahooQuotesRow :exec
INSERT INTO yahoo_quotes (
        created,
        symbol,
        short_name,
        long_name,
        regular_market_price
    )
VALUES (?, ?, ?, ?, ?);
-- name: SaveTipranksRow :exec
INSERT INTO tipranks (
        created,
        symbol,
        news_sentiment,
        analyst_consensus,
        analyst_price_target,
        best_analyst_consensus,
        best_analyst_price_target
    )
VALUES (?, ?, ?, ?, ?, ?, ?);