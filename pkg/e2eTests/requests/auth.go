package requests

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
)

func (c *Client) Login(a rest.LoginRequest) (*http.Response, error) {
	account, _ := json.Marshal(a)
	return http.Post(c.Host+"/auth/login", "application/json", bytes.NewBuffer(account))
}
