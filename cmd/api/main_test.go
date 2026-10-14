package main

import "testing"

func TestInventoryStartupRequiresAuthConfiguration(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("AUTH_JWT_SECRET", "")
	if err := run(); err == nil {
		t.Fatal("started without signing key")
	}
}
