package e2etests_test

import (
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"testing"
)

func TestLaunchServer(t *testing.T) {
	go rest.StartServer(":9000")

	_, err := http.Get("http://localhost:9000")

	if err != nil {
		t.Errorf("Failed to reach server with error: %v", err)
	}
}
