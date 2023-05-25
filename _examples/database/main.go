package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	dbDriver := datastore.DriverSqlite // other drivers are: datastore.DriverPostgresql
	dbURL := "example.sqlite.database"
	dbPoolMax := 1
	dbPrintQueriesToStdout := true

	// connecting to dabase via helper function returns a Bun ORM Type *bun.DB
	db := datastore.NewDBConnection(dbDriver, dbURL, dbPoolMax, dbPrintQueriesToStdout)

	if err := db.Ping(); err != nil {
		fmt.Println("error pinging database" + err.Error())
	}

	// to learn to make queries consult: https://github.com/uptrace/bun

	a := datastore.NewDBHelper(db)
	defer a.Close() // closes the db

	ctx := context.TODO()

	// Migration create ONE OR MORE tables ONLY when they dont exists.
	{

		err := a.Migrate(ctx, (*Book)(nil), (*Dictionary)(nil))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// update multiple
	{
		books := []Book{
			{Id: "book1", Title: "hello"},
			{Id: "book2", Title: "hello world"},
		}
		err := a.UpsertByPKey(ctx, &books)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("---------------")
	// FindByPKey gets ONE record by primary-key (set in struct).
	{
		book := Book{Id: "2"}
		err := a.FindByPKey(ctx, &book)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// FindByColumn gets ONE record via supplied column-name & column-value.
	{
		book := Book{}
		err := a.FindByColumn(ctx, &book, "Id", "book99")
		if err != nil {
			fmt.Println(errors.Is(err, sql.ErrNoRows))
			fmt.Println(err.Error())
		}

		fmt.Println(book)
	}

	fmt.Println("---------------")
	// UpdateByPKey updates ONE OR MORE record by their primary-key (set in struct)
	{
		book := Book{Id: "book1", Title: "hello"}
		err := a.UpdateByPKey(ctx, &book)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// update multiple
	{
		books := []Book{
			{Id: "book1", Title: "hello"},
			{Id: "book2", Title: "hello world"},
		}
		err := a.UpdateByPKey(ctx, &books)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Transaction example: using database transaction
	// Just add the method 'NewWithTx(tx)

	{
		seedBooks := []Book{
			{Id: "book1", Title: "hello"},
			{Id: "book2", Title: "hello world"},
		}
		seedDictionaries := []Dictionary{
			{Id: "oxford1", Title: "Oxford"},
			{Id: "thesarus2", Title: "Thesarus"},
		}

		err := a.Transactional(context.TODO(),

			func(ctx context.Context, tx bun.Tx) error {
				err := a.NewWithTx(tx).Migrate(ctx, (*Book)(nil), (*Dictionary)(nil))
				if err != nil {
					return err
				}

				err = a.NewWithTx(tx).Insert(ctx, &seedBooks, true)
				if err != nil {
					return err
				}

				return a.NewWithTx(tx).Insert(ctx, &seedDictionaries, true)
			},
		)
		if err != nil {
			fmt.Println("error migrating and seeding database" + err.Error())
		}
	}
}
