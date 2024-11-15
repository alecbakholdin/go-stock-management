package zacks

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"stock-management/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZacksDaily(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test", r.FormValue(FormKey))
		w.Write([]byte(`Symbol,Company,Price,Shares,$Chg,%Chg,Industry Rank,Zacks Rank,Value Score,Growth Score,Momentum Score,VGM Score
"AA","Alcoa","41.21","0","0.62","1.53","227","3","C","B","A","B"
"KNDI","Kandi Technologies Group","1.28","1","0.03","2.40","56","NA","NA","NA","NA","NA"
`))
	}))

	s := saver{}
	z := NewDaily(&s, server.URL, "test")
	rows, err := z.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedFetch := []zacksCsvRow{
		{"AA", "Alcoa", 41.21, 0.62, 1.53, 227, 3, "C", "B", "A", "B"},
		{"KNDI", "Kandi Technologies Group", 1.28, 0.03, 2.40, 56, 0, "NA", "NA", "NA", "NA"},
	}
	assert.ElementsMatch(t, expectedFetch, rows)
	if !assert.NoError(t, z.Save(rows)) {
		t.FailNow()
	}

	expectedSave := []models.SaveZacksDailyParams{
		{
			Symbol:        "AA",
			Company:       sql.NullString{String: "Alcoa", Valid: true},
			Price:         41.21,
			DollarChange:  0.62,
			PercentChange: 1.53,
			IndustryRank:  sql.NullInt32{Int32: 227, Valid: true},
			ZacksRank:     sql.NullInt32{Int32: 3, Valid: true},
			ValueScore:    sql.NullString{String: "C", Valid: true},
			GrowthScore:   sql.NullString{String: "B", Valid: true},
			MomentumScore: sql.NullString{String: "A", Valid: true},
			VgmScore:      sql.NullString{String: "B", Valid: true},
		},
		{
			Symbol:        "KNDI",
			Company:       sql.NullString{String: "Kandi Technologies Group", Valid: true},
			Price:         1.28,
			DollarChange:  0.03,
			PercentChange: 2.40,
			IndustryRank:  sql.NullInt32{Int32: 56, Valid: true},
			ZacksRank:     sql.NullInt32{},
			ValueScore:    sql.NullString{},
			GrowthScore:   sql.NullString{},
			MomentumScore: sql.NullString{},
			VgmScore:      sql.NullString{},
		},
	}
	assert.ElementsMatch(t, expectedSave, s.savedRows)
}

type saver struct {
	savedRows []models.SaveZacksDailyParams
}

func (s *saver) SaveZacksDaily(ctx context.Context, arg []models.SaveZacksDailyParams) (int64, error) {
	s.savedRows = arg
	return int64(len(arg)), nil
}