package store

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schemaDDL string

func Migrate(db *sql.DB) error {
	var version int
	err := db.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		return err
	}

	if version == 0 {
		_, err = db.Exec(schemaDDL)
		if err != nil {
			return err
		}
		_, err = db.Exec("PRAGMA user_version = 1")
		if err != nil {
			return err
		}
	}
	return nil
}
