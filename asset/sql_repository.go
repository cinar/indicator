// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

// SQLRepository provides a SQL backed storage facility for financial market data.
type SQLRepository struct {
	// Logger is the slog logger instance.
	Logger *slog.Logger

	// db is the database connection.
	db *sql.DB

	// dialect is the database dialect to use.
	dialect SQLRepositoryDialect

	// assetsQuery is the prepared assets query.
	assetsQuery *sql.Stmt

	// getSinceQuery is the prepared get since query.
	getSinceQuery *sql.Stmt

	// lastDateQuery is the prepared last date query.
	lastDateQuery *sql.Stmt

	// appendQuery is the prepared append query.
	appendQuery *sql.Stmt
}

// NewSQLRepository takes a database driver, URL, and dialect for the asset repository and connects to it.
func NewSQLRepository(dbDriver, dbURL string, dialect SQLRepositoryDialect) (*SQLRepository, error) {
	db, err := sql.Open(dbDriver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect database: %w", err)
	}

	_, err = db.Exec(dialect.CreateTable())
	if err != nil {
		return nil, helper.CloseDatabaseWithError(db, fmt.Errorf("unable to create table: %w", err))
	}

	assetQuery, err := db.Prepare(dialect.Assets())
	if err != nil {
		return nil, helper.CloseDatabaseWithError(db, fmt.Errorf("unable to prepare assets: %w", err))
	}

	getSinceQuery, err := db.Prepare(dialect.GetSince())
	if err != nil {
		return nil, helper.CloseDatabaseWithError(db, fmt.Errorf("unable to prepare get since query: %w", err))
	}

	lastDateQuery, err := db.Prepare(dialect.LastDate())
	if err != nil {
		return nil, helper.CloseDatabaseWithError(db, fmt.Errorf("unable to prepare last date query: %w", err))
	}

	appendQuery, err := db.Prepare(dialect.Append())
	if err != nil {
		return nil, helper.CloseDatabaseWithError(db, fmt.Errorf("unable to prepare append: %w", err))
	}

	repository := &SQLRepository{
		Logger:        slog.Default(),
		db:            db,
		dialect:       dialect,
		assetsQuery:   assetQuery,
		getSinceQuery: getSinceQuery,
		lastDateQuery: lastDateQuery,
		appendQuery:   appendQuery,
	}

	return repository, nil
}

// Close closes the database connection.
func (s *SQLRepository) Close() error {
	return helper.CloseDatabaseWithError(s.db, nil)
}

// Assets returns the names of all assets in the respository.
func (s *SQLRepository) Assets() ([]string, error) {
	rows, err := s.assetsQuery.Query()
	if err != nil {
		return nil, fmt.Errorf("unable to get assets: %w", err)
	}

	defer helper.CloseDatabaseRows(rows)

	var assets []string

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("unable to scan assets: %w", err)
		}

		assets = append(assets, name)
	}

	return assets, nil
}

// Get attempts to return a channel of snapshots for the asset with the given name.
//
// By design, this only returns snapshots from 2000-01-01 onward, unlike
// InMemoryRepository and FileSystemRepository, whose Get returns full
// history with no floor date.
func (s *SQLRepository) Get(name string) (<-chan *Snapshot, error) {
	return s.GetSince(name, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// GetSince attempts to return a channel of snapshots for the asset with the given name since the given date.
//
// The query runs in a background goroutine that owns the underlying database
// rows and sends each scanned snapshot on the returned channel. Callers must
// drain the returned channel to completion (or otherwise ensure it keeps
// being read) so the goroutine can finish and its deferred cleanup can
// release the rows/connection; abandoning the channel partway through leaks
// the goroutine and the underlying database resources, since there is no
// cancellation signal on this interface. See Repository.GetSince for the
// broader interface-level contract.
func (s *SQLRepository) GetSince(name string, date time.Time) (<-chan *Snapshot, error) {
	rows, err := s.getSinceQuery.Query(name, date)
	if err != nil {
		return nil, fmt.Errorf("unable to get since: %w", err)
	}

	snapshots := make(chan *Snapshot)

	go func() {
		defer helper.CloseDatabaseRows(rows)
		defer close(snapshots)

		for rows.Next() {
			snapshot := &Snapshot{}

			err := rows.Scan(
				&snapshot.Date,
				&snapshot.Open,
				&snapshot.High,
				&snapshot.Low,
				&snapshot.Close,
				&snapshot.Volume,
			)
			if err != nil {
				s.Logger.Error("Unable to scan row.", "asset", name, "error", err)
				continue
			}

			snapshots <- snapshot
		}

		if err := rows.Err(); err != nil {
			s.Logger.Error("Unable to iterate rows.", "asset", name, "error", err)
		}
	}()

	return snapshots, nil
}

// LastDate returns the date of the last snapshot for the asset with the given name.
func (s *SQLRepository) LastDate(name string) (time.Time, error) {
	row := s.lastDateQuery.QueryRow(name)

	var date time.Time

	err := row.Scan(&date)
	if err != nil {
		if err == sql.ErrNoRows {
			return date, ErrRepositoryAssetNotFound
		}

		return date, fmt.Errorf("unable to get the last date: %w", err)
	}

	return date, nil
}

// Append adds the given snapshots to the asset with the given name.
func (s *SQLRepository) Append(name string, snapshots <-chan *Snapshot) error {
	var appendErrors []error

	for snapshot := range snapshots {
		_, err := s.appendQuery.Exec(
			name,
			snapshot.Date,
			snapshot.Open,
			snapshot.High,
			snapshot.Low,
			snapshot.Close,
			snapshot.Volume,
		)

		if err != nil {
			appendErrors = append(appendErrors, fmt.Errorf("unable to append snapshot: %w", err))
		}
	}

	return errors.Join(appendErrors...)
}

// Drop drops the snapshots table.
func (s *SQLRepository) Drop() error {
	_, err := s.db.Exec(s.dialect.DropTable())
	if err != nil {
		return fmt.Errorf("unable to drop repository: %w", err)
	}

	return nil
}
