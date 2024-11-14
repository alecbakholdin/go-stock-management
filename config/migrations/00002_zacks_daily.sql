-- +goose Up
-- +goose StatementBegin
CREATE TABLE zacks_daily (
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol VARCHAR(20) NOT NULL,
    company VARCHAR(255),
    price FLOAT NOT NULL,
    dollar_change FLOAT NOT NULL,
    percent_change FLOAT NOT NULL,
    industry_rank INT,
    zacks_rank INT,
    value_score VARCHAR(1),
    growth_score VARCHAR(1),
    momentum_score VARCHAR(1),
    vgm_score VARCHAR(1),

    PRIMARY KEY(created, symbol),
    FOREIGN KEY (symbol) REFERENCES company(symbol)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE zacks_daily;
-- +goose StatementEnd