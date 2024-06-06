package files

import (
	"os/user"
	"path/filepath"
)

// Location
var u, _ = user.Current()
var h = u.HomeDir

var CacheDir, _ = filepath.Abs(filepath.Join(h, ".cache", "aquamarine"))
var Config, _ = filepath.Abs(filepath.Join(h, ".config", "aquamarine", "config.yml"))
var IndexCacheLoc = filepath.Join(CacheDir, "index.json")

func ConstructCachePath(id string) string {
	return filepath.Join(CacheDir, id + ".json")
}
