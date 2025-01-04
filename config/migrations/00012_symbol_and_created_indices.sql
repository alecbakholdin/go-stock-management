-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_yahoo_insights_symbol ON yahoo_insights (symbol);
CREATE INDEX idx_yahoo_insights_created ON yahoo_insights (created);

CREATE INDEX idx_yahoo_quotes_symbol ON yahoo_quotes (symbol);
CREATE INDEX idx_yahoo_quotes_created ON yahoo_quotes (created);

CREATE INDEX idx_zacks_daily_symbol ON zacks_daily (symbol);
CREATE INDEX idx_zacks_daily_created ON zacks_daily (created);

CREATE INDEX idx_zacks_growth_symbol ON zacks_growth (symbol);
CREATE INDEX idx_zacks_growth_created ON zacks_growth (created);

CREATE INDEX idx_tipranks_symbol ON tipranks (symbol);
CREATE INDEX idx_tipranks_created ON tipranks (created);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE yahoo_insights DROP INDEX idx_yahoo_insights_symbol;
ALTER TABLE yahoo_insights DROP INDEX idx_yahoo_insights_created;

ALTER TABLE yahoo_quotes DROP INDEX idx_yahoo_quotes_symbol;
ALTER TABLE yahoo_quotes DROP INDEX idx_yahoo_quotes_created;

ALTER TABLE zacks_daily DROP INDEX idx_zacks_daily_symbol;
ALTER TABLE zacks_daily DROP INDEX idx_zacks_daily_created;

ALTER TABLE zacks_growth DROP INDEX idx_zacks_growth_symbol;
ALTER TABLE zacks_growth DROP INDEX idx_zacks_growth_created;

ALTER TABLE tipranks DROP INDEX idx_tipranks_symbol;
ALTER TABLE tipranks DROP INDEX idx_tipranks_created;
-- +goose StatementEnd
