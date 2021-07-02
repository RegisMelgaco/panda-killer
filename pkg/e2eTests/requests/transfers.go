package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
)

func (c *Client) CreateTransfer(transfer rest.CreateTransferRequest) (resp *http.Response, err error) {
	body, _ := json.Marshal(transfer)
	return http.Post(c.Host+"/transfers", "application/json", bytes.NewBuffer(body))
}

func (c *Client) ListTransfers(accountID int) (*http.Response, error) {
	return http.Get(c.Host + "/transfers/" + fmt.Sprint(accountID))
}
