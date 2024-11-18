package yahoo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock-management/internal/models"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type YahooQuotesSaver interface {
	ListCompanies(context.Context) ([]string, error)
	SaveYahooQuotesRow(context.Context, models.SaveYahooQuotesRowParams) error
}

type yahooQuotesExecutor struct {
	q               YahooQuotesSaver
	crumbUrl        string
	quotesUrlPrefix string
}

func NewQuotes(q YahooQuotesSaver, crumbUrl, quotesUrlPrefix string) *yahooQuotesExecutor {
	return &yahooQuotesExecutor{
		q:               q,
		crumbUrl:        crumbUrl,
		quotesUrlPrefix: quotesUrlPrefix,
	}
}

func (y *yahooQuotesExecutor) Fetch() ([]yahooQuotesJsonRow, error) {
	req, err := y.buildQuotesRequest()
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("error fetching from quotes url"), err)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error %d from yahoo", res.StatusCode) 
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("error reading from quotes body"), err)
	}
	var jsonResponse yahooQuotesJsonResponse
	err = json.Unmarshal(bytes, &jsonResponse)
	if err != nil {
		return nil, errors.Join(errors.New("error unmarshaling json response"), err)
	}

	return jsonResponse.QuoteResponse.Result, nil
}

func (y *yahooQuotesExecutor) buildQuotesRequest() (*http.Request, error) {
	crumb, cookies, err := y.getCrumbAndCookies()
	if err != nil {
		return nil, err
	} 
	companies, err := y.q.ListCompanies(context.Background())
	if err != nil {
		return nil, errors.Join(errors.New("error fetching companyList"), err)
	}
	url, err := url.Parse(y.quotesUrlPrefix)
	if err != nil {
		return nil, errors.Join(errors.New("error parsing quotesUrlPrefix"), err)
	}

	values := url.Query()
	values.Add("crumb", string(crumb))
	values.Add("symbols", strings.Join(companies, ","))
	url.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, errors.Join(errors.New("error creating http request for getting yahoo quotes"), err)
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return req, nil
}

func (y *yahooQuotesExecutor) getCrumbAndCookies() (string, []*http.Cookie, error) {
	req, err := http.NewRequest(http.MethodGet, y.crumbUrl, nil)
	if err != nil {
		return "", nil, errors.Join(errors.New("error creating http request during getCrumb"), err)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
    req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
    req.Header.Add("accept-language", "en-US,en;q=0.9")
    req.Header.Add("cache-control", "max-age=0")
    req.Header.Add("priority", "u=0, i")
    req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not?A_Brand\";v=\"99\"")
    req.Header.Add("sec-ch-ua-mobile", "?0")
    req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
    req.Header.Add("sec-fetch-dest", "document")
    req.Header.Add("sec-fetch-mode", "navigate")
    req.Header.Add("sec-fetch-site", "none")
    req.Header.Add("sec-fetch-user", "?1")
    req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("referrerPolicy", "strict-origin-when-cross-origin")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, errors.Join(errors.New("error fetching from yahoo crumb url"), err)
	} else if res.StatusCode != 200 {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", nil, errors.Join(fmt.Errorf("error reading from crumb request body (status %d)", res.StatusCode), err)
		}
		return "", nil, fmt.Errorf("http error %d reading from crumb url: %s", res.StatusCode, string(bytes))
	}
	crumb, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, errors.Join(errors.New("error reading from crumb body"), err)
	} else if len(crumb) == 0 {
		return "", nil, errors.New("crumb was empty")
	}
	return string(crumb), res.Cookies(), nil
}

type yahooQuotesJsonResponse struct {
	QuoteResponse struct {
		Result []yahooQuotesJsonRow
	}
}

type yahooQuotesJsonRow struct {
	Symbol             string
	RegularMarketPrice float64
}

func (q *yahooQuotesExecutor) Save(rows []yahooQuotesJsonRow) (int, error) {
	created := time.Now()
	n := 0
	for i, row := range rows {
		sqlRow := models.SaveYahooQuotesRowParams{
			Created:            created,
			Symbol:             row.Symbol,
			RegularMarketPrice: row.RegularMarketPrice,
		}
		if err := q.q.SaveYahooQuotesRow(context.Background(), sqlRow); err != nil {
			log.Warnf("Error saving row %d for symbol %s: %s", i, row.Symbol, err.Error())
		} else {
			n += 1
		}
	}
	return n, nil
}
