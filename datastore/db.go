package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

const (
	// Postgresql driver
	DriverPostgresql = "pg"
	// SQlite driver
	DriverSqlite = sqliteshim.ShimName
)

type (
	OrmDB   = *bun.DB
	OrmDbTx = bun.IDB
)

// NewDBConnection establishes a database connection
func NewDBConnection(dbDriver, dsn string, dbPoolMax int, printQueriesToStdout bool) *bun.DB {
	_dbh, err := sql.Open(dbDriver, dsn)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err.Error())
		os.Exit(0)
	}

	var db *bun.DB
	switch dbDriver {
	case DriverSqlite:
		db = bun.NewDB(_dbh, sqlitedialect.New(), bun.WithDiscardUnknownColumns())
	case DriverPostgresql:
		db = bun.NewDB(_dbh, pgdialect.New(), bun.WithDiscardUnknownColumns())
	default:
		fmt.Printf("unknown db driver: %s", dbDriver)
		os.Exit(0)
	}

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(printQueriesToStdout),
	))

	db.SetMaxOpenConns(dbPoolMax * runtime.GOMAXPROCS(0))
	db.SetMaxIdleConns(dbPoolMax * runtime.GOMAXPROCS(0))

	return db
}
