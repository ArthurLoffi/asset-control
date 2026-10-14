package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestEnvironmentPrefersExportedValues(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, key := range []string{"TEST_URI", "TEST_TIMEOUT", "TEST_FLAG"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	contents := "# Local configuration\nexport TEST_URI=\"mongodb://localhost:27017/?retryWrites=true\"\nTEST_TIMEOUT=3s\nTEST_FLAG=false\n"
	if err := os.WriteFile(".env", []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	environment, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	timeout, err := environment.ReadDuration("TEST_TIMEOUT", time.Second)
	if environment.Read("TEST_URI") != "mongodb://localhost:27017/?retryWrites=true" || timeout != 3*time.Second || err != nil {
		t.Fatal("configuration did not match .env values")
	}
	if _, isExported := os.LookupEnv("TEST_URI"); isExported {
		t.Fatal("loading configuration modified the process environment")
	}
	t.Setenv("TEST_URI", "")
	if _, err := environment.Require("TEST_URI"); err == nil || err.Error() != "TEST_URI is required" {
		t.Fatal("an explicitly empty environment variable must not fall back to .env")
	}
	t.Setenv("TEST_FLAG", "true")
	if isSet, err := environment.ReadBool("TEST_FLAG", false); !isSet || err != nil {
		t.Fatal("exported environment variables must override .env values")
	}
	for _, value := range []string{"invalid", "0s", "-1s"} {
		t.Setenv("TEST_TIMEOUT", value)
		if _, err := environment.ReadDuration("TEST_TIMEOUT", time.Second); err == nil {
			t.Fatalf("accepted duration %q", value)
		}
	}
	t.Setenv("TEST_LIST", " a , b")
	if list := environment.ReadList("TEST_LIST", ""); len(list) != 2 || list[0] != "a" || list[1] != "b" {
		t.Fatalf("%q", list)
	}
}

func TestMalformedDotEnvDoesNotExposeSecrets(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("MONGODB_URI=\"private-password"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), ".env") || strings.Contains(err.Error(), "private-password") {
		t.Fatal("expected a safe error identifying the malformed .env file")
	}
}
