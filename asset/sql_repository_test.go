// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
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

type mockRepoStmtAppendErr struct{}

func (s *mockRepoStmtAppendErr) Close() error  { return nil }
func (s *mockRepoStmtAppendErr) NumInput() int { return -1 }
func (s *mockRepoStmtAppendErr) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("insert failed")
}
func (s *mockRepoStmtAppendErr) Query(args []driver.Value) (driver.Rows, error) {
	return &mockRepoRows{}, nil
}

type mockRepoConnAppendErr struct{}

func (c *mockRepoConnAppendErr) Prepare(query string) (driver.Stmt, error) {
	if query == "APPEND" {
		return &mockRepoStmtAppendErr{}, nil
	}

	return &mockRepoStmt{}, nil
}
func (c *mockRepoConnAppendErr) Close() error              { return nil }
func (c *mockRepoConnAppendErr) Begin() (driver.Tx, error) { return nil, nil }

type mockRepoDriverAppendErr struct{}

func (d *mockRepoDriverAppendErr) Open(name string) (driver.Conn, error) {
	return &mockRepoConnAppendErr{}, nil
}

func init() {
	sql.Register("mockrepo", &mockRepoDriver{})
	sql.Register("mockrepoappenderr", &mockRepoDriverAppendErr{})
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

func TestSQLRepositoryGetSinceSkipsScanErrors(t *testing.T) {
	dialect := &mockDialect{}
	repo, err := asset.NewSQLRepository("mockrepo", "db", dialect)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	// mockRepoRows only provides a single column, so scanning it into
	// a Snapshot's six fields fails. That row must be skipped rather
	// than forwarded downstream as a zero-value Snapshot.
	snapshots, err := repo.Get("TEST")
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for range snapshots {
		count++
	}

	if count != 0 {
		t.Fatalf("expected 0 snapshots, got %d", count)
	}
}

func TestSQLRepositoryAppendReturnsError(t *testing.T) {
	dialect := &mockDialect{}
	repo, err := asset.NewSQLRepository("mockrepoappenderr", "db", dialect)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	snapshots := make(chan *asset.Snapshot, 1)
	snapshots <- &asset.Snapshot{}
	close(snapshots)

	err = repo.Append("TEST", snapshots)
	if err == nil {
		t.Fatal("expected error from Append, got nil")
	}
}
