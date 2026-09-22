package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	URL       string
	Namespace string
	TTL       time.Duration
}

func ConfigFromEnv(database string) (Config, error) {
	values, err := godotenv.Read(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, errors.New("could not load .env: check file permissions and KEY=value syntax")
	}
	lookup := func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return values[key]
	}
	cfg := Config{URL: lookup("REDIS_URL"), Namespace: database, TTL: time.Minute}
	if value := lookup("REDIS_CACHE_TTL"); value != "" {
		cfg.TTL, err = time.ParseDuration(value)
		if err != nil || cfg.TTL <= 0 {
			return Config{}, errors.New("REDIS_CACHE_TTL must be a positive Go duration")
		}
	}
	return cfg, nil
}

type Redis struct {
	client    *redis.Client
	namespace string
	ttl       time.Duration
}

// Connect returns nil when Redis is not configured, preserving local development
// without a cache service.
func Connect(ctx context.Context, cfg Config) (*Redis, error) {
	if cfg.URL == "" {
		return nil, nil
	}
	options, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, errors.New("invalid REDIS_URL")
	}
	client := redis.NewClient(options)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}
	return &Redis{client: client, namespace: cfg.Namespace, ttl: cfg.TTL}, nil
}

func (cache *Redis) Close() error { return cache.client.Close() }

func (cache *Redis) versionKey(resource string) string {
	return "asset-control:" + cache.namespace + ":version:" + resource
}

func (cache *Redis) responseKey(resource, request, version string) string {
	digest := sha256.Sum256([]byte(request))
	return "asset-control:" + cache.namespace + ":response:" + resource + ":" + version + ":" + hex.EncodeToString(digest[:])
}

func (cache *Redis) Load(ctx context.Context, resource, request string) ([]byte, string, bool, error) {
	version, err := cache.client.Get(ctx, cache.versionKey(resource)).Result()
	if errors.Is(err, redis.Nil) {
		version = "0"
	} else if err != nil {
		return nil, "", false, err
	}
	body, err := cache.client.Get(ctx, cache.responseKey(resource, request, version)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, version, false, nil
	}
	if err != nil {
		return nil, "", false, err
	}
	current, err := cache.client.Get(ctx, cache.versionKey(resource)).Result()
	if errors.Is(err, redis.Nil) {
		current = "0"
	} else if err != nil {
		return nil, "", false, err
	}
	if current != version {
		return nil, current, false, nil
	}
	return body, version, true, nil
}

func (cache *Redis) Store(ctx context.Context, resource, request, version string, body []byte) error {
	if version == "" {
		return nil
	}
	return cache.client.Set(ctx, cache.responseKey(resource, request, version), body, cache.ttl).Err()
}

func (cache *Redis) Invalidate(ctx context.Context, resource string) error {
	return cache.client.Incr(ctx, cache.versionKey(resource)).Err()
}
