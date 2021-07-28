package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"

	"local/panda-killer/cmd/config"
)

func RunMigrations(env config.EnvVariablesProvider) error {
	dbUrl, err := env.GetDBUrl()
	if err != nil {
		panic(err)
	}

	migrationsUrl, err := env.GetMigrationsFolderUrl()
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(migrationsUrl, dbUrl)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		return err
	}

	version, _, _ := m.Version()

	log.Infof("Migration at version %v", version)

	return nil
}
