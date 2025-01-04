-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_yahoo_insights_symbol_created ON yahoo_insights (symbol, created);
CREATE INDEX idx_yahoo_quotes_symbol_created ON yahoo_quotes (symbol, created);
CREATE INDEX idx_zacks_daily_symbol_created ON zacks_daily (symbol, created);
CREATE INDEX idx_zacks_growth_symbol_created ON zacks_growth (symbol, created);
CREATE INDEX idx_tipranks_symbol_created ON tipranks (symbol, created);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE yahoo_insights DROP INDEX idx_yahoo_insights_symbol_created;
ALTER TABLE yahoo_quotes DROP INDEX idx_yahoo_quotes_symbol_created;
ALTER TABLE zacks_daily DROP INDEX idx_zacks_daily_symbol_created;
ALTER TABLE zacks_growth DROP INDEX idx_zacks_growth_symbol_created;
ALTER TABLE tipranks DROP INDEX idx_tipranks_symbol_created;
-- +goose StatementEnd
