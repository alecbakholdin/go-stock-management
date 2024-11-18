-- +goose Up
-- +goose StatementBegin
ALTER TABLE yahoo_insights ADD COLUMN company_name VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE yahoo_insights DROP COLUMN company_name;
-- +goose StatementEnd
