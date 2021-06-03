package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/inblack67/url-shortner/shortener"
	"github.com/pkg/errors"
)

type redisRepo struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, pingErr := client.Ping(client.Context()).Result()
	return client, pingErr
}

func NewRedisClientRepo(redisURL string) (shortener.RedirectRepo, error) {
	repo := &redisRepo{}
	redisClient, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repo.NewRedisClientRepo")
	}
	repo.client = redisClient
	return repo, nil
}

func (r *redisRepo) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepo) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(r.client.Context(), key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repo.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrorRedirectNotFound, "repo.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["createdAt"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repo.Redirect.Find")
	}
	redirect.CreatedAt = createdAt
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	return redirect, nil
}

func (r *redisRepo) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"createdAt": redirect.CreatedAt,
	}
	_, err := r.client.HSet(r.client.Context(), key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repo.Redirect.Store")
	}
	return nil
}
