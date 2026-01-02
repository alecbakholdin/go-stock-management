package tipranks

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"stock-management/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTipranks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("stickers") == "AAPL" {
			http.ServeFile(w, r, "./tipranks_test.json")
		} else {
			http.Error(w, "incorrect stickers", http.StatusBadRequest)
		}
	}))
	tipranksTable := &testTipranksTable{}
	tipranks := New(tipranksTable, server.URL+"?stickers=")
	rows, err := tipranks.Fetch()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rows))

	expectedAppl := tipranksJsonRow{
		TradingInformationData: tipRanksTradingInformationData{"AAPL"},
		TipRanksEssentialData: tipRanksEssentialData{
			NewsSentiment:            tipRanksNewsSentimentData{NewsSentiment: tipRanksNewsSentimentId{4}},
			AnalystConsensusData:     tipranksAnalystConsensusData{AnalystConsensus: tipRanksNameContainer{"Buy"}},
			BestAnalystConsensusData: tipranksAnalystConsensusData{AnalystConsensus: tipRanksNameContainer{"Sell"}},
			PriceTargetData:          tipRanksPriceTargetData{PriceTarget: 299.1},
			BestPriceTargetData:      tipRanksPriceTargetData{PriceTarget: 301.94},
		},
	}
	assert.Equal(t, expectedAppl, rows[0])

	n, err := tipranks.Save(rows)
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tipranksTable.written))

	expectedApplSql := models.SaveTipranksRowParams{
		Symbol:                 "AAPL",
		NewsSentiment:          sql.NullInt32{Int32: 4, Valid: true},
		AnalystConsensus:       sql.NullString{String: "buy", Valid: true},
		AnalystPriceTarget:     sql.NullFloat64{Float64: 299.1, Valid: true},
		BestAnalystConsensus:   sql.NullString{String: "sell", Valid: true},
		BestAnalystPriceTarget: sql.NullFloat64{Float64: 301.94, Valid: true},
	}
	assert.NotZero(t, tipranksTable.written[0].Created)
	expectedApplSql.Created = tipranksTable.written[0].Created
	assert.Equal(t, expectedApplSql, tipranksTable.written[0])
}

type testTipranksTable struct {
	written []models.SaveTipranksRowParams
}

func (t *testTipranksTable) ListCompanies(_ context.Context) ([]string, error) {
	return []string{"AAPL"}, nil
}

func (t *testTipranksTable) SaveTipranksRow(_ context.Context, row models.SaveTipranksRowParams) error {
	t.written = append(t.written, row)
	return nil
}
