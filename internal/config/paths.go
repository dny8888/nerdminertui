package config

import (
	"os"
	"strings"
)

// ExpandPath expands a path that starts with "~/" to the active user's absolute home directory.
func ExpandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") || p == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if p == "~" {
			return homeDir, nil
		}
		return homeDir + p[1:], nil
	}
	return p, nil
}
