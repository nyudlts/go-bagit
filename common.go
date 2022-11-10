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

func PathExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return err
	} else {
		return err
	}
}
