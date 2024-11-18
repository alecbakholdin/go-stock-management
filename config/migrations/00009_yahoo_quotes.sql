-- +goose Up
-- +goose StatementBegin
CREATE TABLE yahoo_quotes (
    created TIMESTAMP NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    regular_market_price FLOAT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE yahoo_quotes;
-- +goose StatementEnd
