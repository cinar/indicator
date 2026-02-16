// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

import (
	"database/sql"
	"fmt"
	"log"
)

// CloseDatabaseWithError closes the database after an error.
func CloseDatabaseWithError(db *sql.DB, err error) error {
	closeErr := db.Close()
	if closeErr == nil {
		return err
	}

	closeErr = fmt.Errorf("unable to close database: %w", closeErr)

	if err != nil {
		log.Println(closeErr)
		return err
	}

	return closeErr
}

// CloseDatabaseRows closes the database rows.
func CloseDatabaseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Println(err)
	}
}
