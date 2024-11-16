package zacks

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

func TestZacksGrowth(t *testing.T) {
	inputStr := `Symbol,Company,Price,Shares,Growth Score,YR/YR Q0 Growth,LTG %,Last FY Actual,This FY Est,Next FY Est,Q1 Est,Earnings ESP,Next Report Date
"AA","Alcoa","44.02","0","B","150.00","58.86","-2.27","0.89","2.98","0.77","1.13%","1/15/25"
"BBBYQ","Bed Bath & Beyond","0.00","1","NA","0.00","0.00","-0.98","0.00","0.00","0.00","0.00%","NA"
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.FormValue(FormKey))
		w.Write([]byte(inputStr))
	}))
	growthSaver := &zacksGrowthSaver{}
	growthExecutor := NewGrowth(growthSaver, server.URL, "test")

	csvRows, err := growthExecutor.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expectedCsvRows := []zacksGrowthCsvRow{
		{"AA", "Alcoa", 44.02, "B", 150, 58.86, -2.27, 0.89, 2.98, 0.77, 1.13, "1/15/25"},
		{"BBBYQ", "Bed Bath & Beyond", 0, "NA", 0, 0, -0.98, 0, 0, 0, 0, "NA"},
	}
	assert.ElementsMatch(t, expectedCsvRows, csvRows)

	n, err := growthExecutor.Save(csvRows)
	if !assert.NoError(t, err) {
		return
	}
	expectedSqlRows := []models.SaveZacksGrowthRowParams {
		{
			Symbol: "AA",
			Company: sql.NullString{String: "Alcoa", Valid: true},
			Price: 44.02,
			GrowthScore: sql.NullString{String: "B", Valid: true},
			YearOverYearQ0Growth: 150,
			LongTermGrowthPercent: 58.86,
			LastFinancialYearActual: -2.27,
			ThisFinancialYearEst: 0.89,
			NextFinanicalYearEst: 2.98,
			Q1Est: 0.77,
			EarningsExpectedSurprisePrediction: 1.13,
			NextReportDate: sql.NullTime{Time: mustParseDate("1/15/25"), Valid: true},
		},
		{
			Symbol: "BBBYQ",
			Company: sql.NullString{String: "Bed Bath & Beyond", Valid: true},
			Price: 0,
			GrowthScore: sql.NullString{},
			YearOverYearQ0Growth: 0,
			LongTermGrowthPercent: 0,
			LastFinancialYearActual: -0.98,
			ThisFinancialYearEst: 0,
			NextFinanicalYearEst: 0,
			Q1Est: 0,
			EarningsExpectedSurprisePrediction: 0,
			NextReportDate: sql.NullTime{},
		},
	}
	assert.Equal(t, 2, n)
	assert.Equal(t, 2, len(growthSaver.saved))
	assert.EqualValues(t, expectedSqlRows[0], growthSaver.saved[0])
	assert.EqualValues(t, expectedSqlRows[1], growthSaver.saved[1])
}

func mustParseDate(dateStr string) time.Time {
	if t, err := time.Parse("1/2/06", dateStr); err != nil {
		panic(err)
	} else {
		return t
	}
}

type zacksGrowthSaver struct {
	saved []models.SaveZacksGrowthRowParams
}

func (g *zacksGrowthSaver) SaveZacksGrowthRow(_ context.Context, row models.SaveZacksGrowthRowParams) error {
	g.saved = append(g.saved, row)
	return nil
}
