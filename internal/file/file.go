// Handles file names, directories, and caching logic
package file

import (
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// OS file location
var u, _ = user.Current()
var h = u.HomeDir
var CacheDir, _ = filepath.Abs(filepath.Join(h, ".cache", "aquamarine"))
var Config, _ = filepath.Abs(filepath.Join(h, ".config", "aquamarine", "config.yml"))
var IndexCacheLoc = filepath.Join(CacheDir, "index.json")

// Cache TTL
var cacheCutoff = 7 * 24 * time.Hour // A week

func GetCachePath(id string) string {
	return filepath.Join(CacheDir, id + ".json")
}

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
func Cache(b []byte, f string) error {
	if err := os.WriteFile(f, b, 0o644); err != nil {
		return err
	}
	return nil
}
