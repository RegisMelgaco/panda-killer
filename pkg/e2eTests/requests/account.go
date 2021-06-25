package requests

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity/account"
	"net/http"
)

func (c *Client) CreateAccount(a account.Account) (*http.Response, error) {
	account, _ := json.Marshal(a)
	return http.Post(c.Host+"/accounts/", "application/json", bytes.NewBuffer(account))
}

func (c *Client) ListAccounts() (*http.Response, error) {
	return http.Get(c.Host + "/accounts/")
}
