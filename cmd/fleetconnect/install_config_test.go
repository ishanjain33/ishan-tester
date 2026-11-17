package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fleet-connector/internal/config"
	"fleet-connector/internal/config/local"
)

func TestInstallConfigCopiesValidConfig(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "store-042.yaml")
	if err := writeConfig(genConfigOptions{
		path:        src,
		description: "store-042",
		endpoints:   []config.Endpoint{{Upstream: config.Upstream{URL: "localhost:8080"}}},
		authtoken:   "test-token",
	}); err != nil {
		t.Fatalf("gen-config: %v", err)
	}

	dst := filepath.Join(dir, "installed", "config.yaml")
	if err := runInstallConfig([]string{"--path", dst, src}); err != nil {
		t.Fatalf("install-config: %v", err)
	}
	cfg, err := local.New(dst).Load(context.Background())
	if err != nil {
		t.Fatalf("load installed config: %v", err)
	}
	if cfg.Description != "store-042" {
		t.Errorf("description = %q, want store-042", cfg.Description)
	}
}

func TestInstallConfigRejectsInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(src, []byte("schema_version: 99\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "config.yaml")
	err := runInstallConfig([]string{"--path", dst, src})
	if err == nil || !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("want invalid-config error, got %v", err)
	}
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		t.Errorf("invalid config should not be written, stat err: %v", statErr)
	}
}

func TestInstallConfigMissingSourceHints(t *testing.T) {
	err := runInstallConfig([]string{"--path", filepath.Join(t.TempDir(), "c.yaml"), "/nonexistent/store-042.yaml"})
	if err == nil || !strings.Contains(err.Error(), "hint: ") {
		t.Fatalf("want hint for missing source, got %v", err)
	}
}
