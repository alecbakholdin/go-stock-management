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
-- name: SaveZacksGrowthRow :exec
INSERT INTO zacks_growth (
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: SaveYahooInsightsRow :exec
INSERT INTO yahoo_insights (
        symbol,
        short_term,
        mid_term,
        long_term,
        estimated_return,
        fair_value
    )
VALUES (?, ?, ?, ?, ?, ?);