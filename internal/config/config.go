package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SSHConfig struct {
	User       string
	Host       string
	KeyPath    string
	Passphrase string
}

func LoadConfig(path string) (SSHConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return SSHConfig{}, err
	}
	defer f.Close()

	cfg := SSHConfig{}
	section := ""
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		if section != "" && section != "ssh" {
			continue
		}

		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key := strings.TrimSpace(k)
		rawValue := strings.TrimSpace(v)
		value, err := parseTomlString(rawValue)
		if err != nil {
			return SSHConfig{}, fmt.Errorf("invalid value for %q: %w", key, err)
		}

		switch key {
		case "user":
			cfg.User = value
		case "host":
			cfg.Host = value
		case "key_path":
			cfg.KeyPath = value
		case "passphrase":
			cfg.Passphrase = value
		}
	}

	if err := scanner.Err(); err != nil {
		return SSHConfig{}, err
	}

	if cfg.User == "" || cfg.Host == "" || cfg.KeyPath == "" {
		return SSHConfig{}, fmt.Errorf("config.toml must define user, host and key_path")
	}

	return cfg, nil
}

func parseTomlString(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	if strings.HasPrefix(value, "\"") {
		unquoted, err := strconv.Unquote(value)
		if err != nil {
			return "", err
		}
		return unquoted, nil
	}

	return value, nil
}
