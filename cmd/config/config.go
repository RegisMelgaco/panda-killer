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

func getEnvVariable(variableKey, errorMessage string) (string, error) {
	envVariable, isEnvVariableSet := os.LookupEnv(variableKey)
	if !isEnvVariableSet {
		return "", errors.New(errorMessage)
	}
	return envVariable, nil
}

func GetDBUrl() (string, error) {
	return getEnvVariable(dbUrlEnvKey, dbEnvNotSetMessage)
}

func GetMigrationsFolderUrl() (string, error) {
	return getEnvVariable(migrationsFolderUrlEnvKey, migrationsFolderUrlEnvNotSetMessage)
}

func GetRestApiPort() (string, error) {
	return getEnvVariable(grpcApiPortEnvKey, grpcApiPortEnvNotSetMessage)
}

func GetGRPCApiPort() (string, error) {
	return getEnvVariable(restApiPortEnvKey, restApiPortEnvNotSetMessage)
}

func GetAccessSecret() (string, error) {
	return getEnvVariable(accessSecretEnvKey, accessSecretEnvNotSetMessage)
}

func GetDebugMode() bool {
	envVariable, isEnvVariableSet := os.LookupEnv(debugModeEnvKey)
	if !isEnvVariableSet {
		return false
	}
	return strings.ToLower(envVariable) == "true"
}
