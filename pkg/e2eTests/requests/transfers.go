package requests

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
)

func (c *Client) CreateTransfer(authorization string, transfer rest.CreateTransferRequest) (resp *http.Response, err error) {
	body, _ := json.Marshal(transfer)
	req, err := http.NewRequest("POST", c.Host+"/transfers/", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json")
	return (&http.Client{}).Do(req)
}

func (c *Client) ListTransfers(authorization string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.Host+"/transfers/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorization)
	return (&http.Client{}).Do(req)
}
