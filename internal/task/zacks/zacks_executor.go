package zacks

import (
	"errors"
	"net/http"
	"net/url"
	"stock-management/internal/util/csv"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
)

type rowSaver[T any] interface {
	save(time.Time, T) error
}

type zacksExecutor[T interface{Key() string}] struct {
	ms        rowSaver[T]
	url       string
	formValue string
	tableName string
}

func (z *zacksExecutor[TCsv]) Fetch() ([]TCsv, error) {
	values := url.Values{}
	values.Add(FormKey, z.formValue)
	res, err := http.PostForm(z.url, values)
	if err != nil {
		return nil, errors.Join(errors.New("error making POST request"), err)
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New("Got status " + strconv.Itoa(res.StatusCode))
	}

	defer res.Body.Close()
	return csv.Parse(res.Body, new(TCsv))
}

func (z *zacksExecutor[TCsv]) Save(rows []TCsv) (int, error) {
	created := time.Now()
	num := 0
	for i, row := range rows {
		if err := z.ms.save(created, row); err != nil {
			log.Warnf("Zacks %s: error saving row %d for sticker %s: %s", z.tableName, i, row.Key(), err.Error())
			log.Warn("Error saving Zacks ", z.tableName, " row for ", i, ": ", err)
		} else {
			num += 1
		}
	}

	return num, nil
}
