package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"

	"local/panda-killer/cmd/config"
)

func RunMigrations() {
	dbUrl, err := config.GetDBUrl()
	if err != nil {
		panic(err)
	}

	migrationsUrl, err := config.GetMigrationsFolderUrl()
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(migrationsUrl, dbUrl)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		panic(err)
	}

	version, _, _ := m.Version()

	log.Infof("Migration at version %v", version)
}
