package time_series

import (
	"net/url"
	"strings"

	"github.com/TinkLabs/go-webservices/fixer"
)

type Client struct {
	B fixer.Backend
}

func Get(startDate, endDate, baseCurrency string, toCurrencies []string) (*fixer.TimeSeriesResp, error) {
	return getC().Get(startDate, endDate, baseCurrency, toCurrencies)
}

func (c Client) Get(startDate, endDate, baseCurrency string, toCurrencies []string) (resp *fixer.TimeSeriesResp, err error) {
	v := url.Values{}
	if startDate != "" {
		v.Add("start_date", startDate)
	}

	if endDate != "" {
		v.Add("end_date", endDate)
	}

	if baseCurrency != "" {
		v.Add("base", baseCurrency)
	}

	symbols := strings.Join(toCurrencies, ",")
	if symbols != "" {
		v.Add("symbols", symbols)
	}

	path := "/timeseries"
	err = c.B.Call("GET", path, &v, nil, &resp)
	return resp, err
}

func getC() Client {
	return Client{fixer.GetBackend(fixer.PublicBackend)}
}