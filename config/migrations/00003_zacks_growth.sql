-- +goose Up
-- +goose StatementBegin
CREATE TABLE zacks_growth (
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol VARCHAR(20) NOT NULL,
    company VARCHAR(255),
    price FLOAT NOT NULL,
    growth_score VARCHAR(1),
    year_over_year_q0_growth FLOAT NOT NULL,
    long_term_growth_percent FLOAT NOT NULL,
    last_financial_year_actual FLOAT NOT NULL,
    this_financial_year_est FLOAT NOT NULL,
    next_finanical_year_est FLOAT NOT NULL,
    q1_est FLOAT NOT NULL,
    earnings_expected_surprise_prediction FLOAT NOT NULL,
    next_report_date TIMESTAMP,

    PRIMARY KEY(created, symbol),
    FOREIGN KEY (symbol) REFERENCES company(symbol)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE zacks_growth;
-- +goose StatementEnd