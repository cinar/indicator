// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cinar/indicator/v2/helper"
)

// SQLRepository provides a SQL backed storage facility for financial market data.
type SQLRepository struct {
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
		db,
		dialect,
		assetQuery,
		getSinceQuery,
		lastDateQuery,
		appendQuery,
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
func (s *SQLRepository) Get(name string) (<-chan *Snapshot, error) {
	return s.GetSince(name, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// GetSince attempts to return a channel of snapshots for the asset with the given name since the given date.
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
				log.Printf("unable to scan row: %v", err)
			}

			snapshots <- snapshot
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
			return date, fmt.Errorf("unable to find asset")
		}

		return date, fmt.Errorf("unable to get the last date: %w", err)
	}

	return date, nil
}

// Append adds the given snapshots to the asset with the given name.
func (s *SQLRepository) Append(name string, snapshots <-chan *Snapshot) error {
	go func() {
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
				log.Printf("unable to append snapshot: %v", err)
			}
		}
	}()

	return nil
}

// Drop drops the snapshots table.
func (s *SQLRepository) Drop() error {
	_, err := s.db.Exec(s.dialect.DropTable())
	if err != nil {
		return fmt.Errorf("unable to drop repository: %w", err)
	}

	return nil
}
