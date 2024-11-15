package zacks

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"stock-management/internal/csv"
	"stock-management/internal/models"
	"strconv"

	"github.com/labstack/gommon/log"
)

type SaveZacksDaily interface {
	SaveZacksDaily(ctx context.Context, arg []models.SaveZacksDailyParams) (int64, error)
}

type dailyUpdate struct {
	q       SaveZacksDaily
	url     string
	initTab string
}

func NewDaily(q SaveZacksDaily, url, initTab string) *dailyUpdate {
	return &dailyUpdate{
		q:       q,
		url:     url,
		initTab: initTab,
	}
}

func (d *dailyUpdate) Fetch() ([]zacksCsvRow, error) {
	values := url.Values{}
	values.Add(FormKey, d.initTab)
	res, err := http.PostForm(d.url, values)
	if err != nil {
		return nil, errors.Join(errors.New("error making POST request"), err)
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New("Got status " + strconv.Itoa(res.StatusCode))
	}

	defer res.Body.Close()
	return csv.Parse(res.Body, &zacksCsvRow{})
}

func (d *dailyUpdate) Save(rows []zacksCsvRow) error {
	sqlRows := make([]models.SaveZacksDailyParams, len(rows))
	for i, row := range rows {
		sqlRows[i] = models.SaveZacksDailyParams{
			Symbol:        row.Symbol,
			Company:       models.NullStringIfMatch(row.Company, "NA"),
			Price:         row.Price,
			DollarChange:  row.DollarChange,
			PercentChange: row.PercentChange,
			IndustryRank:  models.NullInt32IfZero(row.IndustryRank),
			ZacksRank:     models.NullInt32IfZero(row.ZacksRank),
			ValueScore:    models.NullStringIfMatch(row.ValueScore, "NA"),
			GrowthScore:   models.NullStringIfMatch(row.GrowthScore, "NA"),
			MomentumScore: models.NullStringIfMatch(row.MomentumScore, "NA"),
			VgmScore:      models.NullStringIfMatch(row.VGMScore, "NA"),
		}
	}
	num, err := d.q.SaveZacksDaily(context.Background(), sqlRows)
	if err != nil {
		return err
	} else {
		log.Info("Saved ", num, " Zacks daily rows")
	}

	return nil
}

type zacksCsvRow struct {
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
