package main

import (
	"errors"
	"fmt"
	"io/fs"
	"runtime"
)

// withConfigHint appends a next step to the two config file errors an
// operator is most likely to hit: no file at path, or no permission to
// read or write it (the service's config path needs Administrator/root).
func withConfigHint(path string, err error) error {
	if hint := configHint(path, err, runtime.GOOS); hint != "" {
		return fmt.Errorf("%w\nhint: %s", err, hint)
	}
	return err
}

func configHint(path string, err error, goos string) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return fmt.Sprintf("no config at %s. Pass the path to your config.yaml, or create one with gen-config.", path)
	case errors.Is(err, fs.ErrPermission):
		if goos == "windows" {
			return "run PowerShell as Administrator."
		}
		return "re-run with sudo."
	}
	return ""
}
