package yahoo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"stock-management/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYahooQuotes(t *testing.T) {
	testCrumb := "testCrumb"
	crumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "name",
			Value: "value",
		})
		if _, err := w.Write([]byte(testCrumb)); err != nil {
			t.Errorf("Error writin as part of crumbserver: %s", err.Error())
		}
	}))
	quotesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieMatch := func(c *http.Cookie)bool{ 
			return c.Name == "name" && c.Value == "value"
		}
		if !assert.Equal(t, testCrumb, r.URL.Query().Get("crumb")) {
			w.WriteHeader(http.StatusForbidden)
		} else if !assert.Equal(t, "AAPL,MSFT", r.URL.Query().Get("symbols")) {
			w.WriteHeader(http.StatusBadRequest)
		} else if !assert.True(t, slices.ContainsFunc(r.Cookies(), cookieMatch)){
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			http.ServeFile(w, r, "./yahoo_quotes_test.json")
		}
	}))

	quotesSaver := &testQuotesSaver{}
	quotesExecutor := NewQuotes(quotesSaver, crumbServer.URL, quotesServer.URL)

	jsonRows, err := quotesExecutor.Fetch()
	if !assert.NoError(t, err) || !assert.Equal(t, 2, len(jsonRows)) {
		t.FailNow()
	}
	expectedAapl := yahooQuotesJsonRow{"AAPL", 228.9501}
	expectedMsft := yahooQuotesJsonRow{"MSFT", 415.53}
	assert.Equal(t, expectedAapl, jsonRows[0])
	assert.Equal(t, expectedMsft, jsonRows[1])

	n, err := quotesExecutor.Save(jsonRows)
	sqlRows := quotesSaver.written
	if !assert.NoError(t, err) || !assert.Equal(t, 2, len(sqlRows)) || !assert.Equal(t, 2, n) {
		t.FailNow()
	}
	assert.NotZero(t, sqlRows[0].Created)
	assert.Equal(t, sqlRows[0].Created, sqlRows[1].Created)
	expectedAaplSql := models.SaveYahooQuotesRowParams{Created: sqlRows[0].Created, Symbol: "AAPL", RegularMarketPrice: 228.9501}
	expectedMsftSql := models.SaveYahooQuotesRowParams{Created: sqlRows[1].Created, Symbol: "MSFT", RegularMarketPrice: 415.53}
	assert.Equal(t, expectedAaplSql, sqlRows[0])
	assert.Equal(t, expectedMsftSql, sqlRows[1])
}

func TestYahooQuotesFailsGracefullyOnCrumb429(t *testing.T) {
	crumbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, err := w.Write([]byte("too many requests"))
		if err != nil {
			t.Errorf("error writing to response: %s", err.Error())
		}
	}))
	quotesExecutor := NewQuotes(&testQuotesSaver{}, crumbServer.URL, "gibberish")
	_, err := quotesExecutor.Fetch()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "429")
		assert.Contains(t, err.Error(), "")
	}

}

type testQuotesSaver struct {
	written []models.SaveYahooQuotesRowParams
}

func (t *testQuotesSaver) ListCompanies(_ context.Context) ([]string, error) {
	return []string{"AAPL", "MSFT"}, nil
}

func (t *testQuotesSaver) SaveYahooQuotesRow(_ context.Context, row models.SaveYahooQuotesRowParams) error {
	time.Sleep(time.Millisecond)
	t.written = append(t.written, row)
	return nil
}
