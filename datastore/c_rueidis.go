package datastore

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidiscompat"
)

var (
	_         ICache = (*Rueidis)(nil)
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, vPtr any) error
)

type Rueidis struct {
	store  rueidis.Client
	client rueidiscompat.Cmdable
}

// NewRueidis establishes a redis-cache connection via redis/Rueidis library
func NewRueidis(redisURLs []string, password string, disableCache bool) *Rueidis {
	opts := rueidis.ClientOption{
		Username:     "",
		Password:     password,
		ClientName:   "",
		InitAddress:  redisURLs,
		DisableCache: disableCache,
		ShuffleInit:  true,
	}

	client, err := rueidis.NewClient(opts)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &Rueidis{
		store:  client,
		client: rueidiscompat.NewAdapter(client),
	}
}

func (s *Rueidis) Has(ctx context.Context, key string) bool {
	i, _ := s.client.Exists(ctx, key).Result()
	return i > 0
}

func (s *Rueidis) Get(ctx context.Context, key string, dest any) error {
	val, err := s.client.Cache(time.Second).Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return s.Unmarshal(val, dest)
}

func (s *Rueidis) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	bytes, err := s.Marshal(val)
	if err != nil {
		return err
	}
	_, err = s.client.SetNX(ctx, key, bytes, ttl).Result()
	return err
}

func (s *Rueidis) Del(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *Rueidis) Clear(ctx context.Context) error {
	return s.client.FlushAll(ctx).Err()
}

func (s *Rueidis) Close() error {
	s.store.Close()
	return nil
}

func (s *Rueidis) Marshal(v any) ([]byte, error) {
	if Marshal != nil {
		return Marshal(v)
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes(), err
}

func (s *Rueidis) Unmarshal(data []byte, vPtr any) error {
	if Unmarshal != nil {
		return Unmarshal(data, vPtr)
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(vPtr)
}
