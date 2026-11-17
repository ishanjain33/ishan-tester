package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"fleet-connector/internal/config"
	"fleet-connector/internal/config/local"
)

// runInstallConfig validates a config.yaml generated ahead of time (by
// gen-config, usually on a build machine) and writes it to the service's
// config path. The MSI and the Linux install scripts run it when the
// installer is given FLEETCONNECT_INSTALL_CONFIG, so whoever runs the
// install only needs the installer and that file — no wizard, no gen-config
// on the target. An invalid file fails the install instead of leaving the
// service to fall back to the wizard.
func runInstallConfig(args []string) error {
	fs := flag.NewFlagSet("install-config", flag.ContinueOnError)
	path := fs.String("path", "", "where to write the config (default: the service's config path)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: fleetconnect install-config [--path PATH] <config.yaml>")
	}
	src := fs.Arg(0)

	data, err := os.ReadFile(src)
	if err != nil {
		return withConfigHint(src, fmt.Errorf("read config %q: %w", src, err))
	}
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse config %q: %w", src, err)
	}
	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("config %q is invalid: %w", src, err)
	}

	dst := *path
	if dst == "" {
		dst = defaultConfigPath()
	}
	if err := local.Write(dst, cfg); err != nil {
		return withConfigHint(dst, err)
	}
	fmt.Printf("installed %s to %s\n", src, dst)
	return nil
}
