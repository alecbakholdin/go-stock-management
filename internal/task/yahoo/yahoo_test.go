package yahoo

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"stock-management/internal/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestYahoo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var file *os.File
		var err error
		if file, err = os.Open("./yahoo_test.json"); err != nil {
			panic("error opening file " + err.Error())
		}
		defer file.Close()

		if bytes, err := io.ReadAll(file); err != nil {
			panic("error reading from file: " + err.Error())
		} else {
			w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w.Write(bytes)
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
		{"AAPL", yahooJsonInstrumentInfo{yahooJsonTechnicalEvents{yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bearish"}, yahooJsonOutlook{"Neutral"}}, yahooJsonValuation{"Overvalued", "-6%"}}},
		{"MSFT", yahooJsonInstrumentInfo{yahooJsonTechnicalEvents{yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bullish"}, yahooJsonOutlook{"Bullish"}}, yahooJsonValuation{"Overvalued", "-1%"}}},
	}
	assert.Equal(t, len(expectedJsonRows), len(jsonRows))
	assert.EqualValues(t, expectedJsonRows[0], jsonRows[0])
	assert.EqualValues(t, expectedJsonRows[1], jsonRows[1])

	if !assert.NoError(t, yahooExecutor.Save(jsonRows)) {
		t.FailNow()
	}
	expectedSqlRows := []models.SaveYahooInsightsRowParams{
		{
			Symbol:          "AAPL",
			ShortTerm:       sql.NullString{String: "Bullish", Valid: true},
			MidTerm:         sql.NullString{String: "Bearish", Valid: true},
			LongTerm:        sql.NullString{String: "Neutral", Valid: true},
			EstimatedReturn: sql.NullInt32{Int32: -6, Valid: true},
			FairValue:       sql.NullString{String: "Overvalued", Valid: true},
		},
		{
			Symbol:          "MSFT",
			ShortTerm:       sql.NullString{String: "Bullish", Valid: true},
			MidTerm:         sql.NullString{String: "Bullish", Valid: true},
			LongTerm:        sql.NullString{String: "Bullish", Valid: true},
			EstimatedReturn: sql.NullInt32{Int32: -1, Valid: true},
			FairValue:       sql.NullString{String: "Overvalued", Valid: true},
		},
	}
	assert.Equal(t, len(expectedSqlRows), len(yahooSaver.written))
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
	return nil
}
