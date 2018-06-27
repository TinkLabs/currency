package convert

import (
	"fmt"
	"net/url"

	"github.com/TinkLabs/go-webservices/fixer"
)

type Client struct {
	B fixer.Backend
}

func Convert(from, to string, amount float64) (*fixer.ConvertResp, error) {
	return getC().Convert(from, to, amount)
}

func (c Client) Convert(from, to string, amount float64) (resp *fixer.ConvertResp, err error) {
	v := url.Values{}
	
	v.Add("from", from)
	v.Add("to", to)
	v.Add("amount", fmt.Sprintf("%+v", amount))

	path := "/convert"
	err = c.B.Call("GET", path, &v, nil, &resp)
	return resp, err
}

func getC() Client {
	return Client{fixer.GetBackend(fixer.PublicBackend)}
}
