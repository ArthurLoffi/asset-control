package database

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestConfigFromEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	t.Setenv("MONGODB_DATABASE", "asset_control")
	t.Setenv("MONGODB_TIMEOUT", "")
	cfg, err := ConfigFromEnv()
	if err != nil || cfg.Timeout != 10*time.Second {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
	for _, value := range []string{"invalid", "0s", "-1s"} {
		t.Setenv("MONGODB_TIMEOUT", value)
		if _, err := ConfigFromEnv(); err == nil {
			t.Fatalf("accepted timeout %q", value)
		}
	}
	t.Setenv("MONGODB_TIMEOUT", "2s")
	t.Setenv("MONGODB_DATABASE", " ")
	if _, err := ConfigFromEnv(); err == nil {
		t.Fatal("accepted blank database")
	}
	t.Setenv("MONGODB_DATABASE", "asset_control")
	t.Setenv("MONGODB_URI", "")
	if _, err := ConfigFromEnv(); err == nil {
		t.Fatal("accepted missing URI")
	}
}

func TestConfigLoadsDotEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, key := range []string{"MONGODB_URI", "MONGODB_DATABASE", "MONGODB_TIMEOUT"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	contents := "# Local configuration\nexport MONGODB_URI=\"mongodb://localhost:27017/?retryWrites=true&appName=asset-control\"\nMONGODB_DATABASE='asset_control'\nMONGODB_TIMEOUT=3s\n"
	if err := os.WriteFile(".env", []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.URI != "mongodb://localhost:27017/?retryWrites=true&appName=asset-control" || cfg.Database != "asset_control" || cfg.Timeout != 3*time.Second {
		t.Fatal("configuration did not match .env values")
	}
	if _, exists := os.LookupEnv("MONGODB_URI"); exists {
		t.Fatal("loading configuration modified the process environment")
	}

	t.Setenv("MONGODB_URI", "mongodb://deployment:27017")
	t.Setenv("MONGODB_DATABASE", "deployment")
	t.Setenv("MONGODB_TIMEOUT", "5s")
	cfg, err = ConfigFromEnv()
	if err != nil || cfg.URI != "mongodb://deployment:27017" || cfg.Database != "deployment" || cfg.Timeout != 5*time.Second {
		t.Fatal("exported environment variables must override .env values")
	}

	t.Setenv("MONGODB_URI", "")
	if _, err := ConfigFromEnv(); err == nil || err.Error() != "MONGODB_URI is required" {
		t.Fatal("an explicitly empty environment variable must not fall back to .env")
	}
}

func TestMalformedDotEnvDoesNotExposeSecrets(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("MONGODB_URI=\"private-password"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := ConfigFromEnv()
	if err == nil || !strings.Contains(err.Error(), ".env") || strings.Contains(err.Error(), "private-password") {
		t.Fatal("expected a safe error identifying the malformed .env file")
	}
}
