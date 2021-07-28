package e2etest_test

import (
	"local/panda-killer/pkg/gateway/repository"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	repository.StartPostgresTestContainer()

	code := m.Run()

	repository.FinishPostgresTestContainer()

	os.Exit(code)
}
