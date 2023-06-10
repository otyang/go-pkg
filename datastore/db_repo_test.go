package datastore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func connectDB(dsn string) *bun.DB {
	dbDriver := DriverSqlite
	dbURL := dsn //"file::memory:?cache=shared"
	dbPoolMax := 1
	dbPrintQueriesToStdout := true

	return NewDBConnection(dbDriver, dbURL, dbPoolMax, dbPrintQueriesToStdout)
}

func TestNewDBRepository(t *testing.T) {
	db := connectDB("file::memory:?cache=shared")

	actual := NewDBRepository(db)
	expected := &DBRepository{db: db}

	assert.Equalf(t, expected, actual, "NewDBRepository() = expected %+v but got: %+v", expected, actual)
}

type Book struct {
	Id    string `bun:",pk"`
	Title string `bun:",notnull"`
}

type Dictionary struct {
	Id    string `bun:",pk"`
	Title string `bun:",notnull"`
}

func TestDBRepository_Migrate(t *testing.T) {
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	err := crudRepo.Migrate(context.TODO(), (*Book)(nil), (*Dictionary)(nil))
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)

	_, err = db.NewDropTable().Model(&Book{}).Exec(context.TODO())
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)
}

func TestDBRepository_Create_Find_FindWhere_And_List(t *testing.T) {
	ctx := context.TODO()

	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
		{Id: "book2", Title: "hello world"},
	}

	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	err := crudRepo.Migrate(context.TODO(), (*Book)(nil))
	if err != nil {
		t.Error(err.Error())
	}

	err = crudRepo.Create(ctx, &seedBooks, true)
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)

	// Find
	book := Book{Id: "book1"}
	err = crudRepo.FindByPK(ctx, &book)
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)
	assert.Equalf(t, seedBooks[0], book, "expected %+v but got: %+v", seedBooks[0], book)

	// FindWhere
	addWhere := func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", "book1")
	}
	err = crudRepo.FindWhere(ctx, &book, addWhere)
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)
	assert.Equalf(t, seedBooks[0], book, "expected %+v but got: %+v", seedBooks[0], book)

	// List
	var updBooks []Book
	if err := crudRepo.List(ctx, &updBooks); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	assert.Equalf(t, seedBooks, updBooks, "expected %+v but got: %+v", seedBooks, updBooks)
}

func TestDBRepository_Update(t *testing.T) {
	ctx := context.TODO()
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	if err := crudRepo.Migrate(context.TODO(), (*Book)(nil)); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	// create
	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
	}
	if err := crudRepo.Create(ctx, &seedBooks, true); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	// update
	want := Book{Id: "book1", Title: "hello ==--updated--=="}

	if err := crudRepo.Update(ctx, &want); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	got := Book{Id: "book1"}
	if err := crudRepo.FindByPK(ctx, &got); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	assert.Equalf(t, want, got, "expected %+v but got: %+v", want, got)

	// update bulk
	updatedBooks := []Book{
		{Id: "book1", Title: "hello --updated--"},
		{Id: "book2", Title: "hello world --updated--"},
	}

	if err := crudRepo.UpdateBulk(ctx, &updatedBooks); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	var updBooks []Book
	if err := crudRepo.List(ctx, &updBooks); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	assert.Equalf(t, updatedBooks, updBooks, "expected %+v but got: %+v", updatedBooks, updBooks)
}

func TestDBRepository_Upsert(t *testing.T) {
	ctx := context.TODO()
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	if err := crudRepo.Migrate(context.TODO(), (*Book)(nil)); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	// create
	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
	}
	if err := crudRepo.Create(ctx, &seedBooks, true); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	// upsert
	upsertedBooks := []Book{
		{Id: "book1", Title: "hello --upserted--"},
		{Id: "book2", Title: "hello world"},
	}

	if err := crudRepo.Upsert(ctx, &upsertedBooks); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	var uBooks []Book
	if err := crudRepo.List(ctx, &uBooks); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	assert.Equalf(t, upsertedBooks, uBooks, "expected %+v but got: %+v", upsertedBooks, uBooks)
}

func TestDBRepository_Delete_And_DeleteWhere(t *testing.T) {
	ctx := context.TODO()
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
		{Id: "book2", Title: "hello world"},
		{Id: "book3", Title: "hello world"},
		{Id: "book4", Title: "hello world"},
	}

	err := crudRepo.Migrate(context.TODO(), (*Book)(nil))
	if err != nil {
		t.Error(err.Error())
	}

	if err := crudRepo.Create(ctx, &seedBooks, true); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}

	// DeleteByPK One
	if err := crudRepo.DeleteByPK(ctx, &[]Book{{Id: "book1"}}); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	err = crudRepo.FindByPK(ctx, &Book{Id: "book1"})
	assert.Equalf(t, sql.ErrNoRows, err, "expected %+v but got: %+v", sql.ErrNoRows, err)

	// DeleteByPK Multi
	if err := crudRepo.DeleteByPK(ctx, &[]Book{{Id: "book2"}, {Id: "book3"}}); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	err = crudRepo.FindByPK(ctx, &Book{Id: "book2"})
	assert.Equalf(t, sql.ErrNoRows, err, "expected %+v but got: %+v", sql.ErrNoRows, err)
	err = crudRepo.FindByPK(ctx, &Book{Id: "book3"})
	assert.Equalf(t, sql.ErrNoRows, err, "expected %+v but got: %+v", sql.ErrNoRows, err)

	// DeleteWhere
	addWhereDelete := func(q *bun.DeleteQuery) *bun.DeleteQuery {
		return q.Where("id = ?", "book1")
	}

	if err := crudRepo.DeleteWhere(ctx, &[]Book{{Id: "book1"}}, addWhereDelete); err != nil {
		t.Errorf("expected %+v but got: %+v", nil, err)
	}
	err = crudRepo.FindByPK(ctx, &Book{Id: "book1"})
	assert.Equalf(t, sql.ErrNoRows, err, "expected %+v but got: %+v", sql.ErrNoRows, err)
}

func TestDBRepository_NewWithTx(t *testing.T) {
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	tx, err := db.Begin()
	if err != nil {
		t.Error(err.Error())
	}

	got := crudRepo.NewWithTx(tx)
	var want IDBRepository = &DBRepository{db: tx}

	assert.Equalf(t, want, got, "got this: %+v, but want: %+v", got, want)
}

func TestDBRepository_Transactional(t *testing.T) {
	db := connectDB("file::memory:?cache=shared")
	crudRepo := NewDBRepository(db)

	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
		{Id: "book2", Title: "hello world"},
		{Id: "book3", Title: "hello world"},
		{Id: "book4", Title: "hello world"},
	}

	err := crudRepo.Transactional(context.TODO(),
		func(ctx context.Context, tx bun.Tx) error {
			err := crudRepo.NewWithTx(tx).Migrate(ctx, (*Book)(nil))
			if err != nil {
				return err
			}

			return crudRepo.NewWithTx(tx).Create(ctx, &seedBooks, true)
		},
	)
	assert.Equalf(t, nil, err, "expected %+v but got: %+v", nil, err)
}
