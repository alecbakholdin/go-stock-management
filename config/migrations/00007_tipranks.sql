-- +goose Up
-- +goose StatementBegin
CREATE TABLE tipranks (
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol VARCHAR(20) NOT NULL,
    news_sentiment INT,
    analyst_consensus VARCHAR(20),
    analyst_price_target FLOAT,
    best_analyst_consensus VARCHAR(20),
    best_analyst_price_target FLOAT,
    PRIMARY KEY(created, symbol),
    FOREIGN KEY (symbol) REFERENCES company(symbol)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tipranks;
-- +goose StatementEnd