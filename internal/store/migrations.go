package store

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var SchemaDDL string

func Migrate(db *sql.DB) error {
	_, err := db.Exec(SchemaDDL)
	return err
}
