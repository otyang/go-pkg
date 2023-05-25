package datastore

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	pgc "github.com/patrickmn/go-cache"
)

var (
	_                  ICache = (*GoCache)(nil)
	ErrNotFoundGoCache        = errors.New("item not found")
)

type GoCache struct {
	client *pgc.Cache
}

func NewGoCache() *GoCache {
	defaultExpiration := 8 * time.Second
	cleanupInterval := 1 * time.Minute

	c := pgc.New(defaultExpiration, cleanupInterval)
	return &GoCache{
		client: c,
	}
}

func (s *GoCache) Has(ctx context.Context, key string) bool {
	_, found := s.client.Get(key)
	return found
}

func (s *GoCache) Get(ctx context.Context, key string, dest any) error {
	raw, found := s.client.Get(key)
	if !found {
		return ErrNotFoundGoCache
	}

	v, ok := raw.([]byte)
	if !ok {
		return errors.New("stored value expected to be of type bytes. something else stored")
	}
	return s.Unmarshal(v, dest)
}

func (s *GoCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	bytes, err := s.Marshal(val)
	if err != nil {
		return err
	}
	s.client.Set(key, bytes, ttl)
	return nil
}

func (s *GoCache) Del(ctx context.Context, key string) error {
	s.client.Delete(key)
	return nil
}

func (s *GoCache) Clear(ctx context.Context) error {
	s.client.Flush()
	return nil
}

func (s *GoCache) Close() error {
	return nil
}

func (s *GoCache) Marshal(v any) ([]byte, error) {
	if Marshal != nil {
		return Marshal(v)
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes(), err
}

func (s *GoCache) Unmarshal(data []byte, vPtr any) error {
	if Unmarshal != nil {
		return Unmarshal(data, vPtr)
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(vPtr)
}
