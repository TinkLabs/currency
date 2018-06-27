package fixer

import "fmt"

type ErrorResp struct {
	Success bool `json:"success"`
	Err     Error `json:"error"`
}

type Error struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("code: %d, error_type: %s, info: %s", e.Err.Code, e.Err.Type, e.Err.Info)
}