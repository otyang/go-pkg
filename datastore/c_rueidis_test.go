package datastore

import (
	"context"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

func TestNewRueidis(t *testing.T) {
	type OK struct {
		D string
	}
	ctx := context.Background()
	cache := NewRueidis([]string{"127.0.0.1:6379"}, "", false)

	// set
	_ = cache.Set(ctx, "key3", OK{D: "tracy"}, 5*time.Second)
	_ = cache.Set(ctx, "key4", OK{D: "tracy"}, 5*time.Second)

	var actual OK
	expected := OK{D: "tracy"}

	// get test
	{
		err := cache.Get(ctx, "key3", &actual)
		assert.Equalf(t, expected, actual, "expected %+v but got: %+v", expected, actual)
		assert.Equalf(t, nil, err, "expected nil but got: %s", err)
	}

	// getting item after its been deleted
	{
		cache.Del(ctx, "key4")
		err := cache.Get(ctx, "key4", &actual)

		assert.Equal(t, rueidis.Nil, err, "should return an error")
	}

	// getting item after expiry
	time.Sleep(9 * time.Second)
	{
		expects2 := OK{}
		var actual2 OK
		err := cache.Get(ctx, "key1", &actual2)
		assert.Equalf(t, expects2, actual2, "expected nil but got %#v", actual2)
		assert.Equalf(t, rueidis.Nil, err, "expected nil but got %s", err)
	}
}
