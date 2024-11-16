-- +goose Up
-- +goose StatementBegin
CREATE TABLE company (
    symbol VARCHAR(20) NOT NULL,
    PRIMARY KEY (symbol)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE company;
-- +goose StatementEnd
