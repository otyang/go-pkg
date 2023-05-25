
# Go-pkg

 

>  
>"If I have seen further than others, it is by standing upon the shoulders of giants"  ~ Isaac Newton  
>  

This is not a framework. We're not here to reinvent the wheel for the nth time. 
We believe there are so many great tested libraries out there, why write your own 
when you are not doing something different. This is a pluggable abstraction layer that 
provides a unified interface on top of multiple existing libraries.

So get this use the interface with any library of your choice for interoperability.

# Features

- Config files(json, toml, yaml & env) handling
- Datastore handling (database, events & cache) via any library of your choice
- Logging (via logrus, slog)  via any library of your choice
- Response structure (Http, json)
- Pagination (cursor and offset)
- Validation (either goplay or gookit). you can swap to your library

 


# Requirements
Supports the two most recents versions of Go


# Installation
 


# Usage / Quick Start

## Config
>This does not come with interface purely functional.
```golang
package main

import (
	"fmt"

	"github.com/otyang/go-starter/pkg/config"
)

func main() {

	type ConfigSampleStruct struct {
		Name string `env:"APP_NAME" json:"name" toml:"name" env-default:"Auth"`
	}

	cfgStruct := &ConfigSampleStruct{}      // config struct (must be a pointer)
	pathToConfigFile := "file/location.env" // file could be .env or .json or .toml or .yaml

	/* MustLoad */
	config.MustLoad(pathToConfigFile, cfgStruct) // loads the config and panics on error
	fmt.Println(cfgStruct.Name)

	/* Load */
	err := config.Load(pathToConfigFile, cfgStruct) // loads the config and returns error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfgStruct.Name)
}
```
Config package uses: [ilyakaznacheev/cleanenv](github.com/ilyakaznacheev/cleanenv) for parsing .yaml, .toml, .json, .env files

## Response

>This does not come with interface purely functional it's simple this way.
```golang
package main

import (
	"net/http"
	"github.com/otyang/go-starter/pkg/response"
)

func main() {

	success = true
	httpStatusCode = 400
	msg = "Request was successful"
	errorCode = nil
	data = struct{userName string}{"exampleUser"}
	
	resp := response.New(statusCode, success , msg , &errorCode , data)
	v, _ := json.Marshal(resp)
	fmt.Println(string(v)) 
	/*{
		"success": true,
		"message": "Request was successful",
		"data": {
			"username": "exampleUser",
		}
	}*/


	// A shortcut helper function to create an error response. 
	rsp = response.NewError(http.StatusBadRequest, "this is a bad request", "1000_validation") 
	fmt.Println(rsp)

	// An shortcut helper function to create a success response. 
	rsp = response.NewSuccess(http.StatusOk, "request was successful", data) 
	fmt.Println(rsp)

	// Other helper functions
	rsp = response.UnAuthorised("an error occured on our end", "un_authorized")
	fmt.Println(rsp)



	// Note by default even if the error codes are not defined. response gives a default error code
	// default response code for internal server error is 'internal_Server_Error' so this means
	// out put of 'example1' and 'example2' are same
	example1 := response.InternalServerError("an error occured on our end", "internal_Server_Error")
	example2 := response.InternalServerError("an error occured on our end", "")
	fmt.Println(example1, example2)

	// to disable default errorCode do this
	response.UseDefaultCode = false 
	example3 := response.InternalServerError("an error occured on our end", "")
	fmt.Println(example3) // now the errorCode parameter will be empty
}

// Returning Response as an error:
// you can return response as an error, because it satisfy the golang error interface
func ListSearchEnginesEndPoint() error { 
	return response.InternalServerError("error connecting to database", "")
}

```

Other nice tricks for response package: (define errors into variables)
```golang
var	(
	ErrExpiredToken = response.BadRequest("token is expired", "EXPIRED_TOKEN")
	ErrInvalidTokenFormat = response.BadRequest("invalid token format", "TOKEN_FORMAT")
	ErrAuthenticationRequired = response.BadRequest("you have to login", "AUTH_REQUIRED")
)

overide := ErrInvalidToken.SetMsg("overide the default token is expired message")

fmt.Println(overide, ErrExpiredToken, ErrInvalidTokenFormat, ErrAuthenticationRequired)
```


 


## Logger 
>It provides a set of defined interface through which your library is implemented. Designed around the new proposed Golang SLOG.
>By default Golang's SLOG library is implemented. You can swap to your own
```golang
package main

import (
	"fmt"
	"os"

	"github.com/otyang/go-starter/pkg/logger"
)

func main() {
	var (
		logLevel             = logger.LevelDebug      // logger.LevelWarn, logger.LevelInfo, logger.LevelError
		displayJsonOrText    = logger.DisplayTypeJson // logger.DisplayTypeText
		addSourceCodeAndLine = false
		writeTo              = os.Stdout
	)

	// logger via New Golang SLOG: golang.org/x/exp/slog
	var l = logger.NewSlogLogger(logLevel, displayJsonOrText, addSourceCodeAndLine, writeTo)
	var log logger.Interface = l


	// Log at info level
	log.Info("Hello world at info level")

	// Log With Attributes. This attaches this attributes to every log
	log.With(
		"appName", "helloWorld micro-service",
		"appAddress", "https://helloworld.com:80",
	)

	log.Info("Hello world at info level") // Now notice the difference between this and first example.

	// Helper func that calls log.With() and adds the binary build info
	logger.WithBuildInfo(log)

	log.Reset() // resets the logger to way it was instantiated removing the attributes

	log.Info("cant connect to cache",) 
	log.Info("cant connect to cache", "requestId", "1") 
	fmt.Println("	-----------		")
	log.Debug("cant connect to cache", ) 
	log.Debug("cant connect to cache", "requestId", "1cdhg") 
	fmt.Println("	-----------		")
	log.Fatal("cant connect to cache", ) 
	log.Fatal("cant connect to cache", "requestId", "1cdhg") 
	fmt.Println("	-----------		")
	log.Warn("cant connect to cache", ) 
	log.Warn("cant connect to cache", "requestId", "1cdhg") 
	fmt.Println("	-----------		")
	log.Error("cant connect to cache",)  
	log.Error("cant connect to cache", "requestId", "1cdhg")  
}
```
Logger implements the new Golang SLOG library: https://golang.org/x/exp/slog , you can swap to yours

Logger package interface definition:
```golang
type Interface interface {
	With(args ...any)  
	Reset()
	Debug(message string, args ...any)  
	Error(message string, args ...any) 
	Fatal(message string, args ...any)  
	Info(message string, args ...any)  
	Warn(message string, args ...any)  
}
```


## Datastore - Cache
>This provides a unified usage API for handling cache via any library of your choice, in a way thats easy to swap.
>By default the eko/cache library is used. you can swap to your own library

```golang
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/otyang/go-starter/pkg/datastore"
)

func main() {

	//redis cache via: 	Redis (rueidis) (github.com/redis/rueidis)
	var cache datastore.ICache = datastore.NewRueidis([]string{"127.0.0.1:6379"}, "", false)

	// memory cache via: 	Memory (go-cache) (github.com/patrickmn/go-cache)
	//var cache datastore.ICache = datastore.NewGoCache()

	// setting a cache
	ctx := context.TODO()

	err := cache.Set(ctx, "key", "any thing as value", 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	// getting value from the cache
	var v string
	err = cache.Get(ctx, "key", &v)
	if err != nil {
		fmt.Println("no value in cache:" + err.Error())
	}
	fmt.Println("value from cache:", v)

	// has the key
	time.Sleep(6 * time.Second)
	hasIt := cache.Has(ctx, "key")
	if hasIt {
		fmt.Println("after 6 seconds: key is still saved")
	}

	time.Sleep((6 * time.Second))
	if hasIt := cache.Has(ctx, "key"); !hasIt {
		fmt.Println("after 12 seconds (6 + 6): key is off")
	}

	// delete
	err = cache.Set(ctx, "key", "any thing as value", 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	var b string
	err = cache.Del(ctx, "key")
	if err != nil {
		fmt.Println(err)
	}

	err = cache.Get(ctx, "key", &b)
	if err != nil {
		fmt.Println("no value in cache:" + err.Error())
	}
	fmt.Println("value from cache:", b)
}
```
The cache interface is as follows:
```golang
type ICache interface {
	// Has checks to see if a key exists in the cache
	Has(ctx context.Context, key string) bool
	// Get gets a key from the cache
	Get(ctx context.Context, key string) (any, error)
	// Set sets a key to the cache
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	// Del deletes a key from the cache
	Del(ctx context.Context, key string) error
	// Clear used to flush/clear the cache
	Clear(ctx context.Context) error
	// Close closes the connection
	Close() error
}
```

## Datastore - Database
>Database functions or methoud of usage are numerous and variadic. a unified interface for everything cant be provided.
>To this end the database just provides helper interface for common task executed against a db

```golang
package main

import (
	"context"
	"fmt"

	"github.com/otyang/go-starter/pkg/datastore"
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
	defer a.Close() //closes the db

	ctx := context.TODO()

	// Migration create ONE OR MORE tables ONLY when they dont exists.
	{

		err := a.Migrate(ctx, (*Book)(nil), (*Dictionary)(nil))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

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

	// FindByPKey gets ONE record by primary-key (set in struct).
	{
		book := Book{Id: "1"}
		err := a.FindByPKey(ctx, &book)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// FindByColumn gets ONE record via supplied column-name & column-value.
	{
		book := Book{}
		err := a.FindByColumn(ctx, &book, "columnName", "columnValue")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(book)
	}

	// Transaction example: using database transaction
	// Just add the method 'NewWithTx(passTheTx),

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
```

# API Reference

### /Config
Config comes with the following methods:
-	`config.Load(pathToConfigFile, cfgStructPointer) error`		// loads the config and returns error
-	`config.MustLoad(pathToConfigFile, cfgStructPointer)`		// Panics on error


### /Response
Response come with the following methods api: 
-	`response.NewError(statusCode int, msg string, errorCode string) *response.Response`
-	`response.NewSuccess(statusCode int, msg string, data any) *response.Response`
-	`response.Error() string` Error implements golang error interface
-	`response.SetMsg(msg string) *response.Response` Use to set or overide an already set msg
-	`response.Ok(msg string, data any) *response.Response`
-	`response.Created(msg string, data any) *response.Response`
-	`response.BadRequest(msg string, errorCode string) *response.Response`
-	`response.Unauthorized(msg string, errorCode string) *response.Response`
-	`response.Forbidden(msg string, errorCode string) *response.Response`
-	`response.NotFound(msg string, errorCode string) *response.Response`
-	`response.Conflict(msg string, errorCode string) *response.Response`
-	`response.InternalServerError(msg string, errorCode string) *response.Response`
-	`response.IsJsonErrorGetDetails(err error) (ok bool, e error)` 
-	`response.JSON(w http.ResponseWriter, rsp *Response, headers map[string]string)` for use with std library
-	`response.AsBytes() ([]byte, error)` 
-	`response.DecodeToResponse(s []byte, vPtr any) error`  

### /Loggers 
Loggers come with the following methods:
-	`NewSlogLogger(level LogLevel, displayType LogDisplayType, AddSource bool, output io.Writer) *SlogLogger` New logger via slog
-	`With(args ...any)` // With attached attributes to every log
-	`Reset()`	// resets the logger
-	`Debug(message string, args ...any)` // Debug log a message at DEBUG level
-	`Error(message string, args ...any)` // Error log a message at ERROR level
-	`Fatal(message string, args ...any)` // Fatal log a message at FATAL level
-	`Info(message string, args ...any)` // Info log a message at INFO level
-	`Warn(message string, args ...any)` // Warn log a message at WARN level
-	`WithBuildInfo(logger Interface, args ...any)` // helper func that uses "With" to attache the cmd build information underneath


### /Datastore - Cache
Datastore provides the following interface methods:
-	`Has(ctx context.Context, key string) bool`			Check if the key does exist
-	`Get(ctx context.Context, key string) (any, error)`  Get gets a key from the cache 
-	`Set(ctx context.Context, key string, val any, ttl time.Duration) error`  Set sets a key to the cache 
-	`Del(ctx context.Context, key string) error` // Del deletes a key from the cache
-	`Clear(ctx context.Context) error` Clear used to flush/clear the cache 
-	`Close() error` Close closes the connection 



# Running Tests
To run tests, run the following command:

to do test later
 

# License
This project is licensed under the terms of the [MIT license](/LICENSE).
