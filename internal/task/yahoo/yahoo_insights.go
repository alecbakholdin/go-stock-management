package yahoo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"stock-management/internal/models"
	"stock-management/internal/task/httpunmarshal"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type YahooRowSaver interface {
	ListCompanies(ctx context.Context) ([]string, error)
	SaveYahooInsightsRow(ctx context.Context, row models.SaveYahooInsightsRowParams) error
}

type yahooExecutor struct {
	q         YahooRowSaver
	urlPrefix string
}

func NewInsights(q YahooRowSaver, urlPrefix string) *yahooExecutor {
	return &yahooExecutor{
		q:         q,
		urlPrefix: urlPrefix,
	}
}

func (f *yahooExecutor) Fetch() ([]yahooJsonRow, error) {
	companies, err := f.q.ListCompanies(context.Background())
	if err != nil {
		return nil, errors.Join(errors.New("error listing companies"), err)
	}

	jsonRows := []yahooJsonRow{}
	batchSize := 25
	for i := range (len(companies) / batchSize) + 1 {
		batchStart := i * batchSize
		batchEnd := (i + 1) * batchSize
		if batchEnd > len(companies) {
			batchEnd = len(companies)
		}
		batch := companies[batchStart:batchEnd]

		if jsonRowBatch, err := f.fetchBatch(batch); err != nil {
			return nil, errors.Join(fmt.Errorf("error fetching insights for batch %d", i), err)
		} else {
			jsonRows = append(jsonRows, jsonRowBatch...)
		}
	}
	return jsonRows, nil
}

func (f *yahooExecutor) fetchBatch(companies []string) ([]yahooJsonRow, error) {
	if len(companies) == 0 {
		return []yahooJsonRow{}, nil
	}
	urlPrefix, err := url.Parse(f.urlPrefix)
	if err != nil {
		return nil, errors.Join(errors.New("error parsing yahoo insights url prefix"), err)
	}
	values := urlPrefix.Query()
	values.Add("symbols", strings.Join(companies, ","))
	urlPrefix.RawQuery = values.Encode()

	log.Infof("%s", urlPrefix.String())
	jsonResponse := yahooJsonResponse{}
	if err := httpunmarshal.Get(urlPrefix.String(), &jsonResponse); err != nil {
		return nil, errors.Join(errors.New("error making yahoo insights request"), err)
	}
	return jsonResponse.Finance.Result, nil
}

func (f *yahooExecutor) Save(rows []yahooJsonRow) (int, error) {
	created := time.Now()
	num := 0
	for i, row := range rows {
		estimatedReturn, err := strconv.Atoi(strings.TrimSuffix(row.InstrumentInfo.Valuation.Discount, "%"))
		if err != nil {
			log.Warnf("Error converting %s to int for Yahoo Insights executor for row %d, symbol %s: %s", row.InstrumentInfo.Valuation.Discount, i, row.Symbol, err.Error())
		}
		sqlRow := models.SaveYahooInsightsRowParams{
			Created:         created,
			Symbol:          row.Symbol,
			CompanyName:     models.NullStringIfZero(row.Upsell.CompanyName),
			ShortTerm:       models.NullStringIfZero(row.InstrumentInfo.TechnicalEvents.ShortTermOutlook.Direction),
			MidTerm:         models.NullStringIfZero(row.InstrumentInfo.TechnicalEvents.IntermediateTermOutlook.Direction),
			LongTerm:        models.NullStringIfZero(row.InstrumentInfo.TechnicalEvents.LongTermOutlook.Direction),
			FairValue:       models.NullStringIfZero(row.InstrumentInfo.Valuation.Description),
			EstimatedReturn: sql.NullInt32{Int32: int32(estimatedReturn), Valid: err == nil},
		}
		if err := f.q.SaveYahooInsightsRow(context.Background(), sqlRow); err != nil {
			log.Warnf("Error saving Yahoo Insights row on line %d, sticker %s: %s", i, row.Symbol, err.Error())
		} else {
			num += 1
		}
	}
	return num, nil
}

type yahooJsonResponse struct {
	Finance struct {
		Result []yahooJsonRow
	}
}

type yahooJsonRow struct {
	Symbol         string
	InstrumentInfo yahooJsonInstrumentInfo
	Upsell         yahooJsonUpsell
}

type yahooJsonUpsell struct {
	CompanyName string
}

type yahooJsonInstrumentInfo struct {
	TechnicalEvents yahooJsonTechnicalEvents
	Valuation       yahooJsonValuation
}

type yahooJsonTechnicalEvents struct {
	ShortTermOutlook        yahooJsonOutlook
	IntermediateTermOutlook yahooJsonOutlook
	LongTermOutlook         yahooJsonOutlook
}

type yahooJsonOutlook struct {
	Direction string
}

type yahooJsonValuation struct {
	Description string
	Discount    string
}
