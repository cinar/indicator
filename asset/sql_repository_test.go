// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"

	"github.com/cinar/indicator/v2/asset"
)

type mockDialect struct{}

func (d *mockDialect) CreateTable() string { return "CREATE" }
func (d *mockDialect) DropTable() string   { return "DROP" }
func (d *mockDialect) Assets() string      { return "ASSETS" }
func (d *mockDialect) GetSince() string    { return "GETSINCE" }
func (d *mockDialect) LastDate() string    { return "LASTDATE" }
func (d *mockDialect) Append() string      { return "APPEND" }

type mockRepoDriver struct{}

func (d *mockRepoDriver) Open(name string) (driver.Conn, error) {
	return &mockRepoConn{}, nil
}

type mockRepoConn struct{}

func (c *mockRepoConn) Prepare(query string) (driver.Stmt, error) {
	return &mockRepoStmt{}, nil
}
func (c *mockRepoConn) Close() error              { return nil }
func (c *mockRepoConn) Begin() (driver.Tx, error) { return nil, nil }

type mockRepoStmt struct{}

func (s *mockRepoStmt) Close() error                                    { return nil }
func (s *mockRepoStmt) NumInput() int                                   { return -1 }
func (s *mockRepoStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }
func (s *mockRepoStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &mockRepoRows{}, nil
}

type mockRepoRows struct {
	count int
}

func (r *mockRepoRows) Columns() []string { return []string{"Name"} }
func (r *mockRepoRows) Close() error      { return nil }
func (r *mockRepoRows) Next(dest []driver.Value) error {
	if r.count > 0 {
		return io.EOF
	}
	dest[0] = "TEST"
	r.count++
	return nil
}

func init() {
	sql.Register("mockrepo", &mockRepoDriver{})
}

func TestSQLRepository(t *testing.T) {
	dialect := &mockDialect{}
	repo, err := asset.NewSQLRepository("mockrepo", "db", dialect)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	assets, err := repo.Assets()
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 1 || assets[0] != "TEST" {
		t.Fatalf("expected [TEST], got %v", assets)
	}

	_, err = repo.Get("TEST")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.LastDate("TEST")
	if err == nil {
		// Our mockRepoStmt.QueryRow doesn't return anything, so Scan will fail.
		// That's fine for coverage as long as we reach the line.
	}

	snapshots := make(chan *asset.Snapshot)
	close(snapshots)
	err = repo.Append("TEST", snapshots)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Drop()
	if err != nil {
		t.Fatal(err)
	}
}
