package config

import (
	"errors"
	"os"
)

const (
	dbUrlEnvKey        = "DB_URL"
	DBEnvNotSetMessage = "db url not set in environment variable " + dbUrlEnvKey + " not set"

	migrationsFolderUrlEnvKey           = "MIGRATIONS_FOLDER_URL"
	MigrationsFolderUrlEnvNotSetMessage = "migrations folder url in environment variable " + migrationsFolderUrlEnvKey + " not set"
)

func getEnvVariable(variableKey, errorMessage string) (string, error) {
	url, isEnvVariableSet := os.LookupEnv(variableKey)
	if !isEnvVariableSet {
		return "", errors.New(errorMessage)
	}
	return url, nil
}

func GetDBUrl() (string, error) {
	return getEnvVariable(dbUrlEnvKey, DBEnvNotSetMessage)
}

func GetMigrationsFolderUrl() (string, error) {
	return getEnvVariable(migrationsFolderUrlEnvKey, MigrationsFolderUrlEnvNotSetMessage)
}
