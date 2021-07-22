package config

import (
	"errors"
	"os"
	"strings"
)

const (
	dbUrlEnvKey        = "DB_URL"
	dbEnvNotSetMessage = "db url not set in environment variable " + dbUrlEnvKey + " not set"

	migrationsFolderUrlEnvKey           = "MIGRATIONS_FOLDER_URL"
	migrationsFolderUrlEnvNotSetMessage = "migrations folder url in environment variable " + migrationsFolderUrlEnvKey + " not set"

	restApiPortEnvKey           = "REST_API_ADDRESS"
	restApiPortEnvNotSetMessage = "rest api port environment variable (" + restApiPortEnvKey + ") is not set"

	grpcApiPortEnvKey           = "GRPC_API_ADDRESS"
	grpcApiPortEnvNotSetMessage = "grpc api port environment variable (" + grpcApiPortEnvKey + ") is not set"

	accessSecretEnvKey           = "ACCESS_SECRET"
	accessSecretEnvNotSetMessage = "access secret environment variable (" + accessSecretEnvKey + ") is not set"

	debugModeEnvKey = "DEBUG_MODE"
)

type EnvVariablesProvider interface {
	GetDBUrl() (string, error)
	GetMigrationsFolderUrl() (string, error)
	GetRestApiPort() (string, error)
	GetGRPCApiPort() (string, error)
	GetAccessSecret() (string, error)
	GetDebugMode() bool
	SetTestDBUrl(string) EnvVariablesProvider
}

type EnvVariablesProviderImpl struct {
	testDBURl string
}

func getEnvVariable(variableKey, errorMessage string) (string, error) {
	envVariable, isEnvVariableSet := os.LookupEnv(variableKey)
	if !isEnvVariableSet {
		return "", errors.New(errorMessage)
	}
	return envVariable, nil
}

func (m EnvVariablesProviderImpl) GetDBUrl() (string, error) {
	if len(m.testDBURl) > 0 {
		return m.testDBURl, nil
	}
	return getEnvVariable(dbUrlEnvKey, dbEnvNotSetMessage)
}

func (m EnvVariablesProviderImpl) GetMigrationsFolderUrl() (string, error) {
	return getEnvVariable(migrationsFolderUrlEnvKey, migrationsFolderUrlEnvNotSetMessage)
}

func (m EnvVariablesProviderImpl) GetRestApiPort() (string, error) {
	return getEnvVariable(grpcApiPortEnvKey, grpcApiPortEnvNotSetMessage)
}

func (m EnvVariablesProviderImpl) GetGRPCApiPort() (string, error) {
	return getEnvVariable(restApiPortEnvKey, restApiPortEnvNotSetMessage)
}

func (m EnvVariablesProviderImpl) GetAccessSecret() (string, error) {
	return getEnvVariable(accessSecretEnvKey, accessSecretEnvNotSetMessage)
}

func (m EnvVariablesProviderImpl) GetDebugMode() bool {
	envVariable, isEnvVariableSet := os.LookupEnv(debugModeEnvKey)
	if !isEnvVariableSet {
		return false
	}
	return strings.ToLower(envVariable) == "true"
}

func (m EnvVariablesProviderImpl) SetTestDBUrl(url string) EnvVariablesProvider {
	m.testDBURl = url
	return m
}
