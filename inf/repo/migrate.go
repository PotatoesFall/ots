package repo

import (
	"database/sql"
	"fmt"

	"github.com/PotatoesFall/ots/inf/repo/migrations"
)

func migrate(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback() //nolint:errcheck

	row := tx.QueryRow(`SELECT to_regclass('version') IS NOT NULL;`)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		panic(err)
	}
	if !exists {
		_, err := tx.Exec(`CREATE TABLE version (version INTEGER NOT NULL); INSERT INTO version (version) VALUES (0);`)
		if err != nil {
			panic(err)
		}
	}

	row = tx.QueryRow(`SELECT v.version FROM version v;`)
	var version int
	if err := row.Scan(&version); err != nil {
		panic(err)
	}

	for _, migration := range migrations.ForVersion(version) {
		_, err := tx.Exec(string(migration))
		if err != nil {
			fmt.Println(string(migration))
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
