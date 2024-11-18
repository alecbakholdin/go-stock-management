package yahoo

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"stock-management/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYahoo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if symbols := r.URL.Query().Get("symbols"); symbols != "AAPL,MSFT" {
			t.Errorf("Expected AAPL,MSFT but got %s", symbols)
		} else {
			http.ServeFile(w, r, "./yahoo_insights_test.json")
		}
	}))
	url := server.URL + "?symbols="

	yahooSaver := &yahooSaver{}
	yahooExecutor := NewInsights(yahooSaver, url)

	jsonRows, err := yahooExecutor.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedJsonRows := []yahooJsonRow{
		{"AAPL", yahooJsonInstrumentInfo{yahooJsonTechnicalEvents{yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bearish"}, yahooJsonOutlook{"Neutral"}}, yahooJsonValuation{"Overvalued", "-6%"}}, yahooJsonUpsell{"Apple Inc."}},
		{"MSFT", yahooJsonInstrumentInfo{yahooJsonTechnicalEvents{yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bullish"}}, yahooJsonValuation{"Overvalued", "-1%"}}, yahooJsonUpsell{"Microsoft Corporation"}},
	}
	assert.Equal(t, len(expectedJsonRows), len(jsonRows))
	assert.EqualValues(t, expectedJsonRows[0], jsonRows[0])
	assert.EqualValues(t, expectedJsonRows[1], jsonRows[1])

	n, err := yahooExecutor.Save(jsonRows)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expectedSqlRows := []models.SaveYahooInsightsRowParams{
		{
			Symbol:          "AAPL",
			CompanyName:     sql.NullString{String: "Apple Inc.", Valid: true},
			ShortTerm:       sql.NullString{String: "Bullish", Valid: true},
			MidTerm:         sql.NullString{String: "Bearish", Valid: true},
			LongTerm:        sql.NullString{String: "Neutral", Valid: true},
			EstimatedReturn: sql.NullInt32{Int32: -6, Valid: true},
			FairValue:       sql.NullString{String: "Overvalued", Valid: true},
		},
		{
			Symbol:          "MSFT",
			CompanyName:     sql.NullString{String: "Microsoft Corporation", Valid: true},
			ShortTerm:       sql.NullString{String: "Bullish", Valid: true},
			MidTerm:         sql.NullString{String: "Bullish", Valid: true},
			LongTerm:        sql.NullString{String: "Bullish", Valid: true},
			EstimatedReturn: sql.NullInt32{Int32: -1, Valid: true},
			FairValue:       sql.NullString{String: "Overvalued", Valid: true},
		},
	}
	assert.Equal(t, 2, n)
	assert.Equal(t, len(expectedSqlRows), len(yahooSaver.written))
	assert.NotZero(t, yahooSaver.written[0].Created)
	assert.Equal(t, yahooSaver.written[0].Created, yahooSaver.written[1].Created)
	expectedSqlRows[0].Created = yahooSaver.written[0].Created
	expectedSqlRows[1].Created = yahooSaver.written[1].Created
	assert.Equal(t, expectedSqlRows[0], yahooSaver.written[0])
	assert.Equal(t, expectedSqlRows[1], yahooSaver.written[1])
}

type yahooSaver struct {
	written []models.SaveYahooInsightsRowParams
}

func (y *yahooSaver) ListCompanies(ctx context.Context) ([]string, error) {
	return []string{"AAPL", "MSFT"}, nil
}

func (y *yahooSaver) SaveYahooInsightsRow(ctx context.Context, row models.SaveYahooInsightsRowParams) error {
	y.written = append(y.written, row)
	time.Sleep(time.Millisecond)
	return nil
}
