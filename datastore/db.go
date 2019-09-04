package datastore

import (
	"log"
	"sync"

	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}
var DBH modl.SqlExecutor = DB
var connectOnce sync.Once

func Connect() {
	connectOnce.Do(func() {
		var err error
		url := os.GetEnv("DATABASE_URL")
		DB.Dbx, err = sqlx.Open("postgres", url)
		if err != nil {
			log.Fatalf("Error connecting to Postgres DB using PG* env: %s", err)
		}

		// you w0t m8?!
		DB.Db = DB.Dbx.DB
	})
}

var createSql []string

func Create() {
	if err := DB.CreateTablesIfNotExists(); err != nil {
		log.Fatalf("Error creating tables %s", err)
		return
	}

	for _, query := range createSql {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Error executing query %s: %s", query, err)
		}
	}
}

func Drop() {
	DB.DropTables()
}

func transact(dbh modl.SqlExecutor, fn func(dbh modl.SqlExecutor) error) error {
	var sharedTx bool
	tx, sharedTx := dbh.(*modl.Transaction)
	if !sharedTx {
		var err error
		tx, err = dbh.(*modl.DbMap).Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	if err := fn(tx); err != nil {
		return err
	}

	if !sharedTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
