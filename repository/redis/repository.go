package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type redisReporitory struct {
	client *redis.Client
}

func (r *redisReporitory) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisReporitory) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, err
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

func (r *redisReporitory) Store(*shortener.Redirect) error {
	return nil
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repository := &redisReporitory{}

	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, err
	}

	repository.client = client

	return repository, nil
}
