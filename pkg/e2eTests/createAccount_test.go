package e2etest

import (
	"bytes"
	"encoding/json"
	"local/panda-killer/pkg/domain/entity"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	ts := httptest.NewServer(rest.CreateRouter())
	defer ts.Close()

	account, _ := json.Marshal(entity.Account{})

	resp, _ := http.Post(ts.URL+"/accounts", "application/json", bytes.NewBuffer(account))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to create account with response: %v", resp)
	}
}
