package requests

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
)

func (c *Client) CreateTransfer(transfer rest.CreateTransferRequest) (resp *http.Response, err error) {
	body, _ := json.Marshal(transfer)
	return http.Post(c.Host+"/transfers", "application/json", bytes.NewBuffer(body))
}

func (c *Client) ListTransfers(authorization string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.Host+"/transfers/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorization)
	return (&http.Client{}).Do(req)
}
