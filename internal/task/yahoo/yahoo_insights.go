package yahoo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock-management/internal/models"
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

	urlPrefix, err := url.Parse(f.urlPrefix)
	if err != nil {
		return nil, errors.Join(errors.New("error parsing yahoo insights url prefix"), err)
	}
	values := urlPrefix.Query()
	values.Add("symbols", strings.Join(companies, ","))
	urlPrefix.RawQuery = values.Encode()
	log.Info(urlPrefix)
	
	jsonResponse := yahooJsonResponse{}
	if res, err := http.Get(urlPrefix.String()); err != nil {
		return nil, errors.Join(errors.New("error making get request to yahoo insights"), err)
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.Join(fmt.Errorf("http error %d on yahoo insights request", res.StatusCode))
	} else if bytes, err := io.ReadAll(res.Body); err != nil {
		return nil, errors.Join(errors.New("error reading from response body"), err)
	} else if err := json.Unmarshal(bytes, &jsonResponse); err != nil {
		return nil, errors.Join(errors.New("error unmarshaling response body"), err)
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
