package database

import (
	"context"

	"asset-control/internal/core/config"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func ConfigFromEnv(environment config.Environment) (Config, error) {
	timeout, err := environment.ReadDuration("MONGODB_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{URI: environment.Read("MONGODB_URI"), Database: environment.Read("MONGODB_DATABASE"), Timeout: timeout}
	return cfg, cfg.Validate()
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.URI) == "" {
		return errors.New("MONGODB_URI is required")
	}
	if strings.TrimSpace(cfg.Database) == "" {
		return errors.New("MONGODB_DATABASE is required")
	}
	if cfg.Timeout <= 0 {
		return errors.New("MONGODB_TIMEOUT must be positive")
	}
	return nil
}

// Connect verifies connectivity before returning the shared connection pool.
func Connect(ctx context.Context, cfg Config) (*mongo.Database, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI).SetTimeout(cfg.Timeout))
	if err != nil {
		return nil, fmt.Errorf("configure MongoDB: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = Disconnect(client)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}
	return client.Database(cfg.Database), nil
}

func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
