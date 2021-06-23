package e2etest_test

import (
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLaunchServer(t *testing.T) {
	ts := httptest.NewServer(rest.CreateRouter())
	defer ts.Close()

	_, err := http.Get(ts.URL)

	assert.Nil(t, err)
}
