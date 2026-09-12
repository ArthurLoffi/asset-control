package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func ConfigFromEnv() (Config, error) {
	// Read local defaults without modifying the process environment. A missing
	// file is normal in deployments that inject environment variables directly.
	values, err := godotenv.Read(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Parser errors can contain credentials from the offending line.
		return Config{}, errors.New("could not load .env: check file permissions and KEY=value syntax")
	}
	lookup := func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return values[key]
	}
	cfg := Config{URI: lookup("MONGODB_URI"), Database: lookup("MONGODB_DATABASE"), Timeout: 10 * time.Second}
	if value := lookup("MONGODB_TIMEOUT"); value != "" {
		var err error
		cfg.Timeout, err = time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid MONGODB_TIMEOUT: %w", err)
		}
	}
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
