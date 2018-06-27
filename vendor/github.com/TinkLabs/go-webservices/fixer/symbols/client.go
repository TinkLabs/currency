package symbols

import (
	"github.com/TinkLabs/go-webservices/fixer"
)

type Client struct {
	B fixer.Backend
}

func List() (*fixer.SymbolsResp, error) {
	return getC().List()
}

func (c Client) List() (resp *fixer.SymbolsResp, err error) {
	path := "/symbols"
	err = c.B.Call("GET", path, nil, nil, &resp)
	return resp, err
}

func getC() Client {
	return Client{fixer.GetBackend(fixer.PublicBackend)}
}
