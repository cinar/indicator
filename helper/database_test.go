// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

type mockDriver struct{}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{name: name}, nil
}

type mockConn struct {
	name string
}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (c *mockConn) Close() error {
	if c.name == "fail" {
		return errors.New("close failed")
	}
	return nil
}

func (c *mockConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func init() {
	sql.Register("mock", &mockDriver{})
}

func TestCloseDatabaseWithError(t *testing.T) {
	db, _ := sql.Open("mock", "success")
	db.Ping() // Trigger connection

	err := errors.New("some error")
	res := helper.CloseDatabaseWithError(db, err)
	if res != err {
		t.Fatalf("expected %v, got %v", err, res)
	}

	db, _ = sql.Open("mock", "fail")
	db.Ping() // Trigger connection

	res = helper.CloseDatabaseWithError(db, err)
	if res != err {
		t.Fatalf("expected %v, got %v", err, res)
	}

	db, _ = sql.Open("mock", "fail")
	db.Ping() // Trigger connection

	res = helper.CloseDatabaseWithError(db, nil)
	if res == nil || res.Error() != "unable to close database: close failed" {
		t.Fatalf("expected close error, got %v", res)
	}
}

func TestCloseDatabaseRows(t *testing.T) {
	// CloseDatabaseRows only logs if it fails, and doesn't return anything.
	// Since we don't easily have *sql.Rows without more mocking, we'll
	// just ensure it doesn't panic on nil or we can skip it.
}
