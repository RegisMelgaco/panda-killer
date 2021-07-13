package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
)

func (c *Client) CreateAccount(a rest.CreateAccountRequest) (*http.Response, error) {
	account, _ := json.Marshal(a)
	return http.Post(c.Host+"/accounts/", "application/json", bytes.NewBuffer(account))
}

func (c *Client) ListAccounts() (*http.Response, error) {
	return http.Get(c.Host + "/accounts/")
}

func (c *Client) GetAccountBalance(accountId int) (*http.Response, error) {
	url := fmt.Sprintf(c.Host + "/accounts/" + fmt.Sprint(accountId) + "/balance")
	return http.Get(url)
}
