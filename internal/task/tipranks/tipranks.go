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
	if err := httpunmarshal.Get(t.url+strings.Join(companies, ","), &jsonObj); err != nil {
		return nil, errors.Join(errors.New("error fetching from tipranks"), err)
	}
	return jsonObj.Data, nil
}

func (t *tipranksExecutor) Save(rows []tipranksJsonRow) (int, error) {
	created := time.Now()
	n := 0
	for i, row := range rows {
		analystConsensus := strings.ToLower(row.TipRanksEssentialData.AnalystConsensusData.AnalystConsensus.Name)
		bestAnalystConsensus := strings.ToLower(row.TipRanksEssentialData.BestAnalystConsensusData.AnalystConsensus.Name)
		sqlRow := models.SaveTipranksRowParams{
			Created:                created,
			Symbol:                 row.TradingInformationData.Symbol,
			NewsSentiment:          sql.NullInt32{Int32: int32(row.TipRanksEssentialData.NewsSentiment.NewsSentiment.Id), Valid: true},
			AnalystConsensus:       models.NullStringIfZero(analystConsensus),
			AnalystPriceTarget:     models.NullFloat64IfZero(float64(row.TipRanksEssentialData.PriceTargetData.PriceTarget)),
			BestAnalystConsensus:   models.NullStringIfZero(bestAnalystConsensus),
			BestAnalystPriceTarget: models.NullFloat64IfZero(float64(row.TipRanksEssentialData.BestPriceTargetData.PriceTarget)),
		}
		if err := t.q.SaveTipranksRow(context.Background(), sqlRow); err != nil {
			log.Warnf("Error saving tipranks row %d for symbol %s: %s", i, row.TradingInformationData.Symbol, err.Error())
		} else {
			n += 1
		}
	}
	return n, nil
}

type tipranksJsonResponse struct {
	Data []tipranksJsonRow `json:"assets"`
}

type tipranksJsonRow struct {
	TradingInformationData tipRanksTradingInformationData `json:"tradingInformationData"`
	TipRanksEssentialData  tipRanksEssentialData          `json:"tipRanksEssentialData"`
}

type tipRanksTradingInformationData struct {
	Symbol string `json:"ticker"`
}

type tipRanksEssentialData struct {
	PriceTargetData          tipRanksPriceTargetData      `json:"priceTargetData"`
	BestPriceTargetData      tipRanksPriceTargetData      `json:"bestPriceTargetData"`
	AnalystConsensusData     tipranksAnalystConsensusData `json:"analystsConsensusData"`
	BestAnalystConsensusData tipranksAnalystConsensusData `json:"bestAnalystsConsensusData"`
	NewsSentiment            tipRanksNewsSentimentData    `json:"newsSentimentData"`
}

type tipRanksPriceTargetData struct {
	PriceTarget float64 `json:"priceTarget"`
}

type tipranksAnalystConsensusData struct {
	AnalystConsensus tipRanksNameContainer `json:"analystConsensus"`
}

type tipRanksNameContainer struct {
	Name string `json:"name"`
}

type tipRanksNewsSentimentData struct {
	NewsSentiment tipRanksNewsSentimentId `json:"newsSentiment"`
}

type tipRanksNewsSentimentId struct {
	Id int `json:"id"`
}
