package redisdb

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Redis Client Wrapper
type Redis struct {
	*redis.Client
}

var Ctx = context.TODO()
var keySet = "reductoKeySet"
var urlBannedSet = "urlBannedSet"

// Creates new redis instance with given address
func New(address string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address, Password: "", DB: db,
	})
	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{client}, nil
}

// Saves a given key into the keySet
func (r *Redis) SaveKey(key string) error {
	_, err := r.SAdd(Ctx, keySet, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Returns a key from the keySet
func (r *Redis) GetKey() (string, error) {
	key, err := r.SPop(Ctx, keySet).Result()
	if err != nil {
		return "", err
	}
	return key, err
}

// Retrieves cardinality of keySet
func (r *Redis) KeyPoolSize() (int64, error) {
	val, err := r.SCard(Ctx, keySet).Result()
	if err != nil {
		return -1, err
	}
	return val, nil
}

// Checks URL in list of banned hosts
func (r *Redis) IsHostBanned(url string) (bool, error) {
	val, err := r.SIsMember(Ctx, urlBannedSet, url).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}
