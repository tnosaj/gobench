package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.otters.xyz/jason.tevnan/gobench/internal"
)

// ExecuteNull contains nothing
type ExecuteNull struct {
	Metrics internal.Metrics
}

// ExecStatement will log statement
func (e ExecuteNull) ExecStatement(statement, label string) error {
	logrus.Infof("executing %q with label %q", statement, label)
	return nil
}

// ExecStatementWithReturnBool will execute a statement 's' and return the resulting Boolean
func (e ExecuteNull) ExecStatementWithReturnBool(statement string) (bool, error) {
	logrus.Infof("executing %q", statement)
	return true, nil
}

// ExecStatementWithReturnInt will execute a statement 's' and return the resulting Int
func (e ExecuteNull) ExecStatementWithReturnInt(statement string) (int, error) {
	logrus.Infof("executing %q", statement)
	return 1337, nil
}

// ExecStatementWithReturnRow will execute a statement 's' and return the resulting Row
func (e ExecuteNull) ExecStatementWithReturnRow(statement, label string) (internal.Row, error) {
	logrus.Infof("executing %q with label %q", statement, label)
	return internal.Row{ID: 1337, K: 1338, C: "", Pad: ""}, nil
}

// ExecDDLStatement will execute a statement 's' as a DDL
func (e ExecuteNull) ExecDDLStatement(statement string) error {
	logrus.Infof("executing %q as DDL", statement)
	return nil
}

// Ping checks if the db is up
func (e ExecuteNull) Ping() error {
	logrus.Info("Ping")
	return nil
}

// GetTableExists returns a query to check if dbName.tableName exists
func (e ExecuteNull) GetTableExists(dbName, tableName string) string {
	return fmt.Sprintf("Fancy query to check if '%s' exists in '%s'", tableName, dbName)
}

// Createable returns a query to check if dbName.tableName exists
func (e ExecuteNull) Createable(dbName, tableName string) string {
	return fmt.Sprintf("Fancy query to create '%s.%s'", tableName, dbName)
}
