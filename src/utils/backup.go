package utils

import (
	"fmt"
	"io"
	"obsimcp/src/config"
	"os"
	"path/filepath"
	"time"
)

/*
!!!
When you use mcp-server to delete or overwrite a file,
it will automatically back it up and store the backed-up file in the location specified by your config
!!!
*/

// backup file
func Backupfile(filePath string) (string, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", fmt.Errorf("file does not exist")
    }
    
    timestamp := time.Now().Format("20250501_110405")
    backupfileName := fmt.Sprintf("%s.bak_%s", filepath.Base(filePath), timestamp)
    
    // ensure dir exist
    backupDir := config.Cfg.Backup.Path
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create backup directory: %v", err)
    }
    
    backupPath := filepath.Join(backupDir, backupfileName)

    // Copy file
    src, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer src.Close()

    dst, err := os.Create(backupPath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        return "", err
    }

    return backupPath, nil
}
