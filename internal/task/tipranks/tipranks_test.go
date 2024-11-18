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
		if r.URL.Query().Get("stickers") == "AAPL,RANDOM" {
			http.ServeFile(w, r, "./tipranks_test.json")
		} else {
			http.Error(w, "incorrect stickers", http.StatusBadRequest)
		}
	}))
	tipranksTable := &testTipranksTable{}
	tipranks := New(tipranksTable, server.URL+"?stickers=")
	rows, err := tipranks.Fetch()
	if !assert.NoError(t, err) || !assert.Equal(t, 2, len(rows)) {
		t.FailNow()
	}

	expectedAppl := tipranksJsonRow{"AAPL", 5, tipranksAnalystConsensus{"buy"}, tipranksAnalystConsensus{"sell"}, 245.35, 246.22}
	expectedRandom := tipranksJsonRow{Symbol: "RANDOM"}
	assert.Equal(t, expectedAppl, rows[0])
	assert.Equal(t, expectedRandom, rows[1])


	n, err := tipranks.Save(rows)
	assert.Equal(t, 2, n)
	if !assert.NoError(t, err) || !assert.Equal(t, 2, len(tipranksTable.written)) {
		t.FailNow()
	}
	expectedApplSql := models.SaveTipranksRowParams{
		Symbol:                 "AAPL",
		NewsSentiment:          sql.NullInt32{Int32: 5, Valid: true},
		AnalystConsensus:       sql.NullString{String: "buy", Valid: true},
		AnalystPriceTarget:     sql.NullFloat64{Float64: 245.35, Valid: true},
		BestAnalystConsensus:   sql.NullString{String: "sell", Valid: true},
		BestAnalystPriceTarget: sql.NullFloat64{Float64: 246.22, Valid: true},
	}
	expectedRandomSql := models.SaveTipranksRowParams{Symbol: "RANDOM", NewsSentiment: sql.NullInt32{Int32: 0, Valid: true}}
	assert.NotZero(t, tipranksTable.written[0].Created)
	assert.Equal(t, tipranksTable.written[0].Created, tipranksTable.written[1].Created)
	expectedApplSql.Created = tipranksTable.written[0].Created
	expectedRandomSql.Created = tipranksTable.written[1].Created
	assert.Equal(t, expectedApplSql, tipranksTable.written[0])
	assert.Equal(t, expectedRandomSql, tipranksTable.written[1])
}

type testTipranksTable struct {
	written []models.SaveTipranksRowParams
}

func (t *testTipranksTable) ListCompanies(_ context.Context) ([]string, error) {
	return []string{"AAPL", "RANDOM"}, nil
}

func (t *testTipranksTable) SaveTipranksRow(_ context.Context, row models.SaveTipranksRowParams) error {
	t.written = append(t.written, row)
	return nil
}
