package datastore

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

var _ IDBHelper = (*DBHelper)(nil)

type DBHelper struct {
	db bun.IDB
}

func NewDBHelper(db *bun.DB) *DBHelper {
	return &DBHelper{db: db}
}

func (helper *DBHelper) Close() error {
	db, ok := helper.db.(*bun.DB)
	if !ok {
		return nil
	}

	return db.Close()
}

// Usage: Migrate(ctx, (*StructModel1)(nil), (*StructModel2)(nil), .....)
func (helper *DBHelper) Migrate(ctx context.Context, modelsPtr ...any) error {
	for _, model := range modelsPtr {
		if _, err := helper.db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			errMsg := "failed creating schema resources: " + err.Error()
			return errors.New(errMsg)
		}
	}
	return nil
}

func (helper *DBHelper) UpdateByPKey(ctx context.Context, modelsPtr any) error {
	_, err := helper.db.NewUpdate().Model(modelsPtr).WherePK().Exec(ctx)
	return err
}

func (helper *DBHelper) UpsertByPKey(ctx context.Context, modelsPtr any) error {
	_, err := helper.db.NewInsert().Model(modelsPtr).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func (helper *DBHelper) Insert(ctx context.Context, modelsPtr any, ignoreDupicates bool) error {
	if ignoreDupicates {
		_, err := helper.db.NewInsert().Model(modelsPtr).Ignore().Exec(ctx)
		return err
	}
	_, err := helper.db.NewInsert().Model(modelsPtr).Exec(ctx)
	return err
}

func (helper *DBHelper) FindByPKey(ctx context.Context, modelsPtr any) error {
	err := helper.db.NewSelect().Model(modelsPtr).Limit(1).Scan(ctx)
	return err
}

func (helper *DBHelper) FindByColumn(
	ctx context.Context, modelsPtr any, columnName string, columnValue any,
) error {
	err := helper.db.
		NewSelect().
		Model(modelsPtr).
		Where("? = ?", bun.Ident(columnName), columnValue).
		Limit(1).
		Scan(ctx)
	return err
}

func (helper *DBHelper) ListAll(ctx context.Context, modelsPtr any) error {
	err := helper.db.NewSelect().Model(modelsPtr).Scan(ctx)
	return err
}

func (helper *DBHelper) ListByColumn(
	ctx context.Context, modelsPtr any, columnName string, columnValue any,
) error {
	err := helper.db.
		NewSelect().
		Model(modelsPtr).
		Where("? = ?", bun.Ident(columnName), columnValue).
		Scan(ctx)
	return err
}

func (helper *DBHelper) DeleteByPKey(ctx context.Context, modelsPtr any) error {
	_, err := helper.db.NewDelete().Model(modelsPtr).WherePK().Exec(ctx)
	return err
}

func (helper *DBHelper) DeleteByColumn(ctx context.Context, modelsPtr, colName string, colValue any) error {
	_, err := helper.db.NewDelete().
		Model(modelsPtr).
		Where(
			"? = ?", bun.Ident(colName), colValue,
		).Exec(ctx)
	return err
}

// NewWithTx returns a clone of DBHelper, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
// as the new dbConnection
func (helper *DBHelper) NewWithTx(tx bun.Tx) IDBHelper {
	return &DBHelper{db: tx}
}

// Transactional simplifies transactions code, by automatically:
// starting a transaction, rolling back the transaction if an error
// is encountered & finally commiting the transaction if no error.
//
//	func Example_transactionUsage(ctx context.Context) {
//		var example IDBHelper = nil
//		example.Transactional(
//			ctx,
//			func(ctx context.Context, tx bun.Tx) error {
//				example.NewWithTx(tx).FindByPKey(ctx, nil)
//				example.NewWithTx(tx).UpdateByPKey(ctx, nil)
//				return nil
//			},
//		)
//	}
func (helper *DBHelper) Transactional(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return helper.db.RunInTx(ctx, nil, fn)
}
