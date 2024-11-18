package zacks

import (
	"context"
	"stock-management/internal/models"
	"time"
)

type DailyRowSaver interface {
	SaveZacksDailyRow(ctx context.Context, arg models.SaveZacksDailyRowParams) error
}

type dailyUpdate struct {
	q DailyRowSaver
}

func NewDaily(q DailyRowSaver, url, formValue string) *zacksExecutor[zacksDailyCSvRow] {
	return &zacksExecutor[zacksDailyCSvRow]{
		ms:        &dailyUpdate{q: q},
		tableName: "Daily",
		formValue: formValue,
		url:       url,
	}
}

func (d *dailyUpdate) save(created time.Time, csvRow zacksDailyCSvRow) error {
	sqlRow := models.SaveZacksDailyRowParams{
		Created:       created,
		Symbol:        csvRow.Symbol,
		Company:       models.NullStringIfMatch(csvRow.Company, "NA"),
		Price:         csvRow.Price,
		DollarChange:  csvRow.DollarChange,
		PercentChange: csvRow.PercentChange,
		IndustryRank:  models.NullInt32IfZero(csvRow.IndustryRank),
		ZacksRank:     models.NullInt32IfZero(csvRow.ZacksRank),
		ValueScore:    models.NullStringIfMatch(csvRow.ValueScore, "NA"),
		GrowthScore:   models.NullStringIfMatch(csvRow.GrowthScore, "NA"),
		MomentumScore: models.NullStringIfMatch(csvRow.MomentumScore, "NA"),
		VgmScore:      models.NullStringIfMatch(csvRow.VGMScore, "NA"),
	}
	return d.q.SaveZacksDailyRow(context.Background(), sqlRow)
}

type zacksDailyCSvRow struct {
	Symbol        string
	Company       string
	Price         float64
	DollarChange  float64 `csv:"$Chg"`
	PercentChange float64 `csv:"%Chg"`
	IndustryRank  int32   `csv:"Industry Rank"`
	ZacksRank     int32   `csv:"Zacks Rank"`
	ValueScore    string  `csv:"Value Score"`
	GrowthScore   string  `csv:"Growth Score"`
	MomentumScore string  `csv:"Momentum Score"`
	VGMScore      string  `csv:"VGM Score"`
}

func (z zacksDailyCSvRow) Key() string {
	return z.Symbol
}
