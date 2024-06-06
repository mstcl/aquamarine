package utils

import (
	"os"
	"time"
)

var cacheCutoff = 7 * 24 * time.Hour // A week

// Return true if cached data needs syncing
func ShouldSync(f string) (bool, error) {
	info, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		} else {
			return false, err
		}
	}

	now := time.Now()
	if diff := now.Sub(info.ModTime()); diff > cacheCutoff {
		return true, nil
	}

	return false, nil
}

// Cache content b to file loc
func CacheFile(b []byte, f string) error {
	if err := os.WriteFile(f, b, 0o644); err != nil {
		return err
	}
	return nil
}
