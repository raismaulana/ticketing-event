package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCase interface {
	Set(key string, value interface{}) error
	Get(key string, v interface{}) (interface{}, error)
}

type redisCase struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisCase(rdb *redis.Client) RedisCase {
	return &redisCase{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (service *redisCase) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(key)
	log.Println(string(json))
	errSet := service.rdb.Set(service.ctx, key, string(json), time.Second*300).Err()
	if errSet != nil {
		log.Println(errSet)
		return errSet
	}

	return nil
}

func (service *redisCase) Get(key string, v interface{}) (interface{}, error) {

	if val, err := service.rdb.Get(service.ctx, key).Result(); err == redis.Nil {
		log.Println("%key does not exist")
		return nil, err
	} else if err != nil {
		log.Println(err)
		return nil, err
	} else {
		errs := json.Unmarshal([]byte(val), &v)
		return v, errs
	}
}

func (service *redisCase) Delete(key string) {
	service.rdb.Del(service.ctx, key)
}
