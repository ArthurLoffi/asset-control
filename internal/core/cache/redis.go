package cache

import (
	"context"

	"asset-control/internal/core/config"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	URL       string
	Namespace string
	TTL       time.Duration
}

func ConfigFromEnv(environment config.Environment, database string) (Config, error) {
	ttl, err := environment.ReadDuration("REDIS_CACHE_TTL", time.Minute)
	return Config{URL: environment.Read("REDIS_URL"), Namespace: database, TTL: ttl}, err
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
