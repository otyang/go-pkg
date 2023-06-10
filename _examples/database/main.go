package main

import (
	"context"
	"fmt"
	"log"

	"github.com/otyang/go-pkg/datastore"
	"github.com/uptrace/bun"
)

type Book struct {
	Id    string `bun:",pk"`
	Title string `bun:",notnull"`
}

type Dictionary struct {
	Id    string `bun:",pk"`
	Title string `bun:",notnull"`
}

func main() {
	ctx := context.TODO()

	// connecting to dabase via helper function returns a Bun ORM Type *bun.DB

	dbDriver := datastore.DriverSqlite // other drivers are: datastore.DriverPostgresql
	dbURL := "file::memory:?cache=shared"
	dbPoolMax := 1
	dbPrintQueriesToStdout := true
	db := datastore.NewDBConnection(dbDriver, dbURL, dbPoolMax, dbPrintQueriesToStdout)

	if err := db.Ping(); err != nil {
		fmt.Println("error pinging database" + err.Error())
	}

	// to learn to make queries consult: https://github.com/uptrace/bun

	// Using the repository helper
	crud := datastore.NewDBRepository(db)
	defer db.Close()

	// migration
	err := crud.Migrate(context.TODO(), (*Book)(nil), (*Dictionary)(nil))
	if err != nil {
		log.Fatal(err.Error())
	}

	// create
	seedBooks := []Book{
		{Id: "book1", Title: "hello"},
		{Id: "book2", Title: "hello world"},
	}
	err = crud.Create(ctx, &seedBooks, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Find
	book := Book{Id: "book1"}
	err = crud.FindByPK(ctx, &book)
	if err != nil {
		log.Fatal(err.Error())
	}

	// FindWhere
	addWhere := func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", "book1")
	}
	err = crud.FindWhere(ctx, &book, addWhere)
	if err != nil {
		log.Fatal(err.Error())
	}

	// List
	var updBooks []Book
	err = crud.List(ctx, &updBooks)
	if err != nil {
		log.Fatal(err.Error())
	}

	// update
	want := Book{Id: "book1", Title: "hello ==--updated--=="}
	err = crud.Update(ctx, &want)
	if err != nil {
		log.Fatal(err.Error())
	}

	// upsert
	upsertedBooks := []Book{
		{Id: "book1", Title: "hello --upserted--"},
		{Id: "book2", Title: "hello world"},
	}
	err = crud.Upsert(ctx, &upsertedBooks)
	if err != nil {
		log.Fatal(err.Error())
	}

	// update bulk
	updatedBooks := []Book{
		{Id: "book1", Title: "hello --updated--"},
		{Id: "book2", Title: "hello world --updated--"},
	}
	err = crud.UpdateBulk(ctx, &updatedBooks)
	if err != nil {
		log.Fatal(err.Error())
	}

	// DeleteByPK One
	err = crud.DeleteByPK(ctx, &[]Book{{Id: "book1"}})
	if err != nil {
		log.Fatal(err.Error())
	}

	// DeleteByPK Multi
	err = crud.DeleteByPK(ctx, &[]Book{{Id: "book1"}, {Id: "book2"}})
	if err != nil {
		log.Fatal(err.Error())
	}

	// DeleteWhere
	addWhereDelete := func(q *bun.DeleteQuery) *bun.DeleteQuery {
		return q.Where("id = ?", "book1")
	}

	if err := crud.DeleteWhere(ctx, &[]Book{{Id: "book1"}}, addWhereDelete); err != nil {
		log.Fatal(err.Error())
	}

	// Transaction example: using database transaction
	// Just add the method 'NewWithTx(tx)
	manyBooks := []Book{
		{Id: "book1", Title: "hello"},
		{Id: "book2", Title: "hello world"},
		{Id: "book3", Title: "hello world"},
		{Id: "book4", Title: "hello world"},
	}

	err = crud.Transactional(context.TODO(),
		func(ctx context.Context, tx bun.Tx) error {
			err := crud.NewWithTx(tx).Migrate(ctx, (*Book)(nil))
			if err != nil {
				return err
			}

			return crud.NewWithTx(tx).Create(ctx, &manyBooks, true)
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
