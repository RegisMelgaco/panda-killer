package requests

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/transfer"
	"net/http"
)

func (c *Client) CreateTransfer(transfer transfer.Transfer) (resp *http.Response, err error) {
	body, _ := json.Marshal(transfer)
	return http.Post(c.Host+"/transfers", "application/json", bytes.NewBuffer(body))
}
