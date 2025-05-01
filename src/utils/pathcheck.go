package utils

import (
	"obsimcp/src/config"
	"os"
	"strings"
)

// Check if a path is under vault
func CheckIllegalPath(path string) bool {
    vaultPath := config.Cfg.Vault.Path

    if !strings.HasPrefix(path, vaultPath) {
        return true
    }

    return false
}

// Check is a markdown file
func CheckIsMd(path string) bool {
    if !strings.HasSuffix(path, ".md") {
        return false
    }

    return true
}

// Check is exist
func CheckIsExist(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        // file not exist
        return false
    }

    return true
}
