package datastore

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

var _ IDBRepository = (*DBRepository)(nil)

type (
	SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery
	DeleteCriteria func(*bun.DeleteQuery) *bun.DeleteQuery
)

type DBRepository struct {
	db bun.IDB
}

func NewDBRepository(db *bun.DB) *DBRepository {
	return &DBRepository{db: db}
}

// Usage: Migrate(ctx, (*StructModel1)(nil), (*StructModel2)(nil), .....)
func (r *DBRepository) Migrate(ctx context.Context, modelsPtr ...any) error {
	for _, model := range modelsPtr {
		if _, err := r.db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			errMsg := "failed creating schema resources: " + err.Error()
			return errors.New(errMsg)
		}
	}
	return nil
}

func (r *DBRepository) Create(ctx context.Context, model any, ignoreDupicates bool) error {
	if ignoreDupicates {
		_, err := r.db.NewInsert().Model(model).Ignore().Returning("*").Exec(ctx)
		return err
	}
	_, err := r.db.NewInsert().Model(model).Returning("*").Exec(ctx)
	return err
}

func (r *DBRepository) Upsert(ctx context.Context, modelsPtr any) error {
	_, err := r.db.NewInsert().Model(modelsPtr).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func (r *DBRepository) FindByPK(ctx context.Context, modelPtr any) error {
	return r.db.NewSelect().Model(modelPtr).WherePK().Limit(1).Scan(ctx)
}

func (r *DBRepository) FindWhere(ctx context.Context, modelPtr any, sc ...SelectCriteria) error {
	q := r.db.NewSelect().Model(modelPtr)

	for i := range sc {
		q.Apply(sc[i])
	}

	return q.Limit(1).Scan(ctx)
}

func (r *DBRepository) List(ctx context.Context, modelPtr any, sc ...SelectCriteria) error {
	q := r.db.NewSelect().Model(modelPtr)

	for i := range sc {
		q.Apply(sc[i])
	}

	return q.Scan(ctx)
}

func (r *DBRepository) Update(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewUpdate().Model(modelPtr).WherePK().Returning("*").Exec(ctx)
	return err
}

func (r *DBRepository) UpdateBulk(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewUpdate().Model(modelPtr).WherePK().Bulk().Returning("*").Exec(ctx)
	return err
}

func (r *DBRepository) DeleteByPK(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewDelete().Model(modelPtr).WherePK().Exec(ctx)
	return err
}

func (r *DBRepository) DeleteWhere(ctx context.Context, modelPtr any, dc ...DeleteCriteria) error {
	q := r.db.NewDelete().Model(modelPtr)

	for i := range dc {
		q.Apply(dc[i])
	}

	return q.Scan(ctx)
}

// NewWithTx returns a clone of Repository, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
// as the new dbConnection
func (r *DBRepository) NewWithTx(tx bun.Tx) IDBRepository {
	return &DBRepository{
		db: tx,
	}
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
func (r *DBRepository) Transactional(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return r.db.RunInTx(ctx, nil, fn)
}
