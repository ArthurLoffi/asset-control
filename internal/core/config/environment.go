package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Environment reads exported variables first and falls back to the optional
// .env file, without modifying the process environment.
type Environment struct{ dotEnv map[string]string }

func Load() (Environment, error) {
	dotEnv, err := godotenv.Read(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Parser errors can contain credentials from the offending line.
		return Environment{}, errors.New("could not load .env: check file permissions and KEY=value syntax")
	}
	return Environment{dotEnv: dotEnv}, nil
}

// Read returns "" when the key is unset; an exported empty value never falls back to .env.
func (e Environment) Read(key string) string {
	if value, isExported := os.LookupEnv(key); isExported {
		return value
	}
	return e.dotEnv[key]
}

func (e Environment) Require(key string) (string, error) {
	value := e.Read(key)
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func (e Environment) ReadDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := e.Read(key)
	if value == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive Go duration", key)
	}
	return duration, nil
}

func (e Environment) ReadBool(key string, fallback bool) (bool, error) {
	value := e.Read(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return parsed, nil
}

// ReadList splits a comma-separated value and trims each entry.
func (e Environment) ReadList(key string, fallback string) []string {
	value := e.Read(key)
	if value == "" {
		value = fallback
	}
	items := strings.Split(value, ",")
	for index := range items {
		items[index] = strings.TrimSpace(items[index])
	}
	return items
}
