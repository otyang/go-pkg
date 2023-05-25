package main

import (
	"context"
	"fmt"
	"time"

	"github.com/otyang/go-pkg/datastore"
)

func main() {
	// redis cache via: 	Redis (rueidis) (github.com/redis/rueidis)
	var cache datastore.ICache = datastore.NewRueidis([]string{"127.0.0.1:6379"}, "", false)

	// memory cache via: 	Memory (go-cache) (github.com/patrickmn/go-cache)
	// var cache datastore.ICache = datastore.NewGoCache()

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
