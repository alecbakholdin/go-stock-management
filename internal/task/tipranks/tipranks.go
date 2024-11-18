package tipranks

import (
	"context"
	"database/sql"
	"errors"
	"stock-management/internal/models"
	"stock-management/internal/task/httpunmarshal"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type TipranksSaver interface {
	ListCompanies(context.Context) ([]string, error)
	SaveTipranksRow(context.Context, models.SaveTipranksRowParams) error
}

type tipranksExecutor struct {
	q   TipranksSaver
	url string
}

func New(q TipranksSaver, url string) *tipranksExecutor {
	return &tipranksExecutor{
		q:   q,
		url: url,
	}
}

func (t *tipranksExecutor) Fetch() ([]tipranksJsonRow, error) {
	companies, err := t.q.ListCompanies(context.Background())
	if err != nil {
		return nil, errors.Join(errors.New("error fetching companies"), err)
	}
	var jsonObj tipranksJsonResponse
	if err := httpunmarshal.Get(t.url + strings.Join(companies, ","), &jsonObj); err != nil {
		return nil, errors.Join(errors.New("error fetching from tipranks"), err)
	}
	return jsonObj.Data, nil
}

func (t *tipranksExecutor) Save(rows []tipranksJsonRow) (int, error) {
	created := time.Now()
	n := 0
	for i, row := range rows {
		sqlRow := models.SaveTipranksRowParams{
			Created:                created,
			Symbol:                 row.Symbol,
			NewsSentiment:          sql.NullInt32{Int32: int32(row.NewsSentiment), Valid: true},
			AnalystConsensus:       models.NullStringIfZero(row.AnalystConsensus.Consensus),
			AnalystPriceTarget:     models.NullFloat64IfZero(float64(row.PriceTarget)),
			BestAnalystConsensus:   models.NullStringIfZero(row.BestAnalystConsensus.Consensus),
			BestAnalystPriceTarget: models.NullFloat64IfZero(float64(row.BestPriceTarget)),
		}
		if err := t.q.SaveTipranksRow(context.Background(), sqlRow); err != nil {
			log.Warnf("Error saving tipranks row %d for symbol %s: %s", i, row.Symbol, err.Error())
		} else {
			n += 1
		}
	}
	return n, nil
}

type tipranksJsonResponse struct {
	Data []tipranksJsonRow
}

type tipranksJsonRow struct {
	Symbol               string `json:"ticker"`
	NewsSentiment        int
	AnalystConsensus     tipranksAnalystConsensus
	BestAnalystConsensus tipranksAnalystConsensus
	PriceTarget          float64
	BestPriceTarget      float64
}

type tipranksAnalystConsensus struct {
	Consensus string
}
