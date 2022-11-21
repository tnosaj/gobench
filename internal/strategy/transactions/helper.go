package transactions

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
)

func (st TransactionReadWrite) automigrate(action string) error {
	logrus.Debug("automatically migrating mysql")
	var driver database.Driver
	var err error
	switch st.S.DB {
	case "mysql":
		driver, err = mysqlmigrate.WithInstance(st.S.DBConnection, &mysqlmigrate.Config{})
		if err != nil {
			logrus.Errorf("Failed to create migration connection %s", err)
			return err
		}

	case "postgres":
		logrus.Debug("automatically migrating postgres")
		driver, err = postgres.WithInstance(st.S.DBConnection, &postgres.Config{})
		if err != nil {
			logrus.Errorf("Failed to create migration connection %s", err)
			return err
		}
	}
	m, _ := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/%s-transactions/", "tmp", st.S.DB),
		st.S.DB,
		driver,
	)

	if err != nil {
		logrus.Errorf("Failed to initialize migration connection %s", err)
		return err
	}
	logrus.Debugf("succesfully migrated %s", st.S.DB)
	switch action {
	case "prepare":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logrus.Errorf("Failed to migrate db: %s", err)
			return err
		}

	case "cleanup":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			logrus.Errorf("Failed to migrate db: %s", err)
			return err
		}
	}
	return nil
}
