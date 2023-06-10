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
type IDBRepository interface {
	// Migration create ONE OR MORE tables ONLY when they dont exists.
	Migrate(ctx context.Context, modelsPtr ...any) error
	// UpdateByPK updates a record by their primary-key (set in struct)
	Update(ctx context.Context, modelsPtr any) error
	// UpdateBulk updates multiple rows via primarykey
	UpdateBulk(ctx context.Context, modelPtr any) error
	// Upsert updates ONE OR MORE record. if the record doesn't exist, it inserts it.
	Upsert(ctx context.Context, modelsPtr any) error
	// Create inserts ONE OR MORE record.
	Create(ctx context.Context, modelPtr any, ignoreDupicates bool) error
	// FindByPK gets ONE record by primary-key (set in struct). [limit 1]
	FindByPK(ctx context.Context, modelPtr any) error
	// FindByColumn gets record(s) via supplied criteria.  [limit 1]
	FindWhere(ctx context.Context, modelPtr any, sc ...SelectCriteria) error
	// List records of a table via criteria. Useful for loading settings from db
	List(ctx context.Context, modelPtr any, sc ...SelectCriteria) error
	// DeleteByPK deletes record(s) using primary key in struct
	DeleteByPK(ctx context.Context, modelsPtr any) error
	// DeleteWhere deletes records(s) via criteria
	DeleteWhere(ctx context.Context, modelsPtr any, dc ...DeleteCriteria) error

	// NewWithTx returns a clone of DBHelper, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
	// as the new dbConnection
	NewWithTx(tx bun.Tx) IDBRepository
	// Transactional simplifies transactions code, by automatically:
	//
	// starting a transaction, rolling back the transaction if an error
	// is encountered & finally commiting the transaction if no error.
	//
	// Transactional(ctx, func(ctx context.Context, tx bun.Tx) error {
	// 		err := NewWithTx(tx).Migrate(ctx, (*Book)(nil))
	// 		if err != nil {
	// 			return err
	// 		}

	// 		return NewWithTx(tx).Create(ctx, &seedBooks, true)
	// 	})
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
