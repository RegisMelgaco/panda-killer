package e2etest_test

import (
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLaunchServer(t *testing.T) {
	ts := httptest.NewServer(rest.CreateRouter())
	defer ts.Close()

	_, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("Failed to reach server with error: %v", err)
	}
}
