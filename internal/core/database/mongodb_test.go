package database

import (
	"testing"
	"time"

	"asset-control/internal/core/config"
)

func TestConfigFromEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	t.Setenv("MONGODB_DATABASE", "asset_control")
	t.Setenv("MONGODB_TIMEOUT", "")
	environment, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := ConfigFromEnv(environment)
	if err != nil || cfg.Timeout != 10*time.Second {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
	for _, value := range []string{"invalid", "0s", "-1s"} {
		t.Setenv("MONGODB_TIMEOUT", value)
		if _, err := ConfigFromEnv(environment); err == nil {
			t.Fatalf("accepted timeout %q", value)
		}
	}
	t.Setenv("MONGODB_TIMEOUT", "2s")
	t.Setenv("MONGODB_DATABASE", " ")
	if _, err := ConfigFromEnv(environment); err == nil {
		t.Fatal("accepted blank database")
	}
	t.Setenv("MONGODB_DATABASE", "asset_control")
	t.Setenv("MONGODB_URI", "")
	if _, err := ConfigFromEnv(environment); err == nil {
		t.Fatal("accepted missing URI")
	}
}
