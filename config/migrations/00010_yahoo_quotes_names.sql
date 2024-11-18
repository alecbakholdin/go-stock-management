-- +goose Up
-- +goose StatementBegin
ALTER TABLE yahoo_quotes ADD short_name VARCHAR(100);
ALTER TABLE yahoo_quotes ADD long_name VARCHAR(255);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE yahoo_quotes DROP COLUMN short_name;
ALTER TABLE yahoo_quotes DROP COLUMN long_name;
-- +goose StatementEnd