package datastore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
)

// IDBHelper is an interface that provides quick helpers for handling database operations
type IDBHelper interface {
	// Close: closes the connection
	Close() error
	// Migration create ONE OR MORE tables ONLY when they dont exists.
	Migrate(ctx context.Context, modelsPtr ...any) error
	// UpdateByPKey updates ONE OR MORE record by their primary-key (set in struct)
	UpdateByPKey(ctx context.Context, modelsPtr any) error
	// UpsertByPKey updates ONE OR MORE record by their primary-key and if the record
	// doesn't exist, it inserts it.
	UpsertByPKey(ctx context.Context, modelsPtr any) error
	// Insert inserts ONE OR MORE record.
	Insert(ctx context.Context, modelsPtr any, ignoreDupicates bool) error
	// FindByPKey gets ONE record by primary-key (set in struct). [limit 1]
	FindByPKey(ctx context.Context, modelsPtr any) error
	// FindByColumn gets a record via supplied column-name & column-value.  [limit 1]
	FindByColumn(ctx context.Context, modelsPtr any, columnName string, columnValue any) error
	// List all records of a table. Useful for loading settings from db
	ListAll(ctx context.Context, modelsPtr any) error
	// List records of a table via supplied column and column value
	ListByColumn(ctx context.Context, modelsPtr any, columnName string, columnValue any) error
	// DeleteByPKey deletes a record using primary key in struct
	DeleteByPKey(ctx context.Context, modelsPtr any) error
	// DeleteByColumn deletes ONE OR MORE record via supplied column-name & column-value
	DeleteByColumn(ctx context.Context, modelsPtr, columnName string, columnValue any) error

	// NewWithTx returns a clone of DBHelper, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
	// as the new dbConnection
	NewWithTx(tx bun.Tx) IDBHelper
	// Transactional simplifies transactions code, by automatically:
	//
	// starting a transaction, rolling back the transaction if an error
	// is encountered & finally commiting the transaction if no error.
	// func Example_transactionUsage(ctx context.Context) {
	// 	var example IDBHelper = nil
	// 	example.Transactional(
	// 		ctx,
	// 		func(ctx context.Context, tx bun.Tx) error {
	// 			example.NewWithTx(tx).FindByPKey(ctx, nil)
	// 			example.NewWithTx(tx).UpdateByPKey(ctx, nil)
	// 			return nil
	// 		},
	// 	)
	// }
	Transactional(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
}

// ICache is an interface that guides & ensure the use of different external cache library, in a way
// thats easy to swap.
type ICache interface {
	// Has checks to see if a key exists in the cache
	Has(ctx context.Context, key string) bool
	// Get gets a key from the cache
	Get(ctx context.Context, key string, dest any) error
	// Set sets a key to the cache
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	// Del deletes a key from the cache
	Del(ctx context.Context, key string) error
	// Clear used to flush/clear the cache
	Clear(ctx context.Context) error
	// Close closes the connection
	Close() error
}

// IsErrNotFound checks if the error returned from any of the datastore library used here
// e.g redis, database etc is an error not found type
func IsErrNotFound(err error) bool {
	switch {
	case
		rueidis.IsRedisNil(err),
		errors.Is(err, sql.ErrNoRows),
		errors.Is(err, ErrNotFoundGoCache):
		return true
	}
	return false
}
