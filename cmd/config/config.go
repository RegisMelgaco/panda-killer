package config

import (
	"errors"
	"os"
)

const (
	dbUrlEnvKey        = "DB_URL"
	dbEnvNotSetMessage = "db url not set in environment variable " + dbUrlEnvKey + " not set"

	migrationsFolderUrlEnvKey           = "MIGRATIONS_FOLDER_URL"
	migrationsFolderUrlEnvNotSetMessage = "migrations folder url in environment variable " + migrationsFolderUrlEnvKey + " not set"

	restApiPortEnvKey           = "REST_API_ADDRESS"
	restApiPortEnvNotSetMessage = "rest api port environment variable (" + restApiPortEnvKey + ") is not set"

	accessSecretEnvKey           = "ACCESS_SECRET"
	accessSecretEnvNotSetMessage = "access secret environment variable (" + accessSecretEnvKey + ") is not set"
)

func getEnvVariable(variableKey, errorMessage string) (string, error) {
	url, isEnvVariableSet := os.LookupEnv(variableKey)
	if !isEnvVariableSet {
		return "", errors.New(errorMessage)
	}
	return url, nil
}

func GetDBUrl() (string, error) {
	return getEnvVariable(dbUrlEnvKey, dbEnvNotSetMessage)
}

func GetMigrationsFolderUrl() (string, error) {
	return getEnvVariable(migrationsFolderUrlEnvKey, migrationsFolderUrlEnvNotSetMessage)
}

func GetRestApiPort() (string, error) {
	return getEnvVariable(restApiPortEnvKey, restApiPortEnvNotSetMessage)
}

func GetAccessSecret() (string, error) {
	return getEnvVariable(accessSecretEnvKey, accessSecretEnvNotSetMessage)
}
