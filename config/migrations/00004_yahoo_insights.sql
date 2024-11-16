-- +goose Up
-- +goose StatementBegin
CREATE TABLE yahoo_insights (
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol VARCHAR(20) NOT NULL,
    short_term VARCHAR(20),
    mid_term VARCHAR(20),
    long_term VARCHAR(20),
    fair_value VARCHAR(20),
    estimated_return INT,
    PRIMARY KEY(created, symbol),
    FOREIGN KEY (symbol) REFERENCES company(symbol)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE yahoo_insights;
-- +goose StatementEnd