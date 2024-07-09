package db

import (
	"context"

	"github.com/sirupsen/logrus"
)

// ExecuteNull contains nothing
type ExecuteNull struct {
	Metrics Metrics
}

// ExecStatement will log statement
func (e ExecuteNull) ExecStatement(statement interface{}, label string) error {
	logrus.Infof("executing %q with label %q", statement, label)
	return nil
}

// ExecStatement will log statement
func (e ExecuteNull) ExecInterfaceStatement(statement interface{}, label string) error {
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
func (e ExecuteNull) ExecStatementWithReturnRow(statement, label string) (interface{}, error) {
	logrus.Infof("executing %q with label %q", statement, label)
	var returnedRow interface{}
	//returnedRow = interface{ID: 1337, K: 1338, C: "", Pad: ""}
	return returnedRow, nil
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

func (e ExecuteNull) AutoMigrateUP(folderName string) error {
	logrus.Infof("MigrateUp from folder '%s'", folderName)
	return nil
}

func (e ExecuteNull) AutoMigrateDown(folderName string) error {
	logrus.Infof("MigrateDown from folder '%s'", folderName)
	return nil
}

func (e ExecuteNull) Shutdown(c context.Context) {
	logrus.Infof("Shutting down nulldb")
}
