package main

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"
)

func TestConfigHint(t *testing.T) {
	notExist := fmt.Errorf("read config: %w", fs.ErrNotExist)
	denied := fmt.Errorf("write config: %w", fs.ErrPermission)

	if got := configHint("/etc/fleetconnect/config.yaml", notExist, "linux"); !strings.Contains(got, "no config at /etc/fleetconnect/config.yaml") {
		t.Errorf("not-exist hint = %q", got)
	}
	if got := configHint("x", denied, "windows"); !strings.Contains(got, "Administrator") {
		t.Errorf("windows permission hint = %q", got)
	}
	if got := configHint("x", denied, "linux"); !strings.Contains(got, "sudo") {
		t.Errorf("linux permission hint = %q", got)
	}
	if got := configHint("x", fmt.Errorf("parse error"), "linux"); got != "" {
		t.Errorf("unrelated error got hint %q", got)
	}
}

func TestWithConfigHintKeepsWrappedError(t *testing.T) {
	err := withConfigHint("x", fmt.Errorf("read: %w", fs.ErrNotExist))
	if !strings.Contains(err.Error(), "hint: ") {
		t.Errorf("missing hint: %v", err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("hint dropped the wrapped error")
	}
}
