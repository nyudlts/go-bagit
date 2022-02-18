package go_bagit

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

var currentTime = time.Now()

func GetSoftwareAgent() string {
	const mod = "github.com/nyudlts/go-bagit"

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Sprintf("go-bagit <https://%s>", mod)
	}

	version := "(unknown)"
	for _, dep := range info.Deps {
		if dep.Path == mod {
			version = dep.Version
			break
		}
	}

	return fmt.Sprintf("go-bagit %s <https://%s>", version, mod)
}

func DirectoryExists(path string) error {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() == true {
			return nil
		}
		return fmt.Errorf("Path is not a Directory")
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Path does not exist")
	} else {
		return fmt.Errorf("Unknown Error")
	}
}
