package httpunmarshal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Get(url string, obj any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(errors.New("error creating http request"), err)
	}
	return Do(req, obj)
}

func Do(req *http.Request, obj any) error {
	if res, err := http.DefaultClient.Do(req); err != nil {
		return errors.Join(errors.New("error getting http response"), err)
	} else if res.StatusCode != http.StatusOK {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Join(fmt.Errorf("error reading body from %d response", res.StatusCode), err)
		}
		return fmt.Errorf("http error %d making http request: %s", res.StatusCode, string(bytes))
	} else if bytes, err := io.ReadAll(res.Body); err != nil {
		return errors.Join(errors.New("error reading from response body"), err)
	} else if err := json.Unmarshal(bytes, obj); err != nil{
		return errors.Join(errors.New("error unmarshaling json object"), err)
	}
	return nil
}
