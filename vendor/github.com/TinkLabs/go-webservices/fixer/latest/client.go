package latest

import (
	"net/url"
	"strings"

	"github.com/TinkLabs/go-webservices/fixer"
)

type Client struct {
	B fixer.Backend
}

func Get(baseCurrency string, toCurrencies []string) (*fixer.LatestResp, error) {
	return getC().Get(baseCurrency, toCurrencies)
}

func (c Client) Get(baseCurrency string, toCurrencies []string) (resp *fixer.LatestResp, err error) {
	v := url.Values{}
	if baseCurrency != "" {
		v.Add("base", baseCurrency)
	}

	symbols := strings.Join(toCurrencies, ",")
	if symbols != "" {
		v.Add("symbols", symbols)
	}

	path := "/latest"
	err = c.B.Call("GET", path, &v, nil, &resp)
	return resp, err
}

func getC() Client {
	return Client{fixer.GetBackend(fixer.PublicBackend)}
}
