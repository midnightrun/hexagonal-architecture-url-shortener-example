package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("RedisRepository Find: %v -> %w", data, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("RedisRepository Find: %v -> %w", data, shortener.ErrRedirectNotFound)
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("RedisRepository Find: %v -> %w", createdAt, err)
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return fmt.Errorf("RedisRepository Store: %v -> %w", data, err)
	}

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
		return nil, fmt.Errorf("RedisRepository newRedisClient: %v -> %w", client, err)
	}

	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repository := &redisRepository{}

	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, fmt.Errorf("RedisRepository NewRedisRepository: %v -> %w", client, err)
	}

	repository.client = client

	return repository, nil
}
