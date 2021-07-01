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
