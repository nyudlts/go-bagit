package go_bagit

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"time"
)

var manifestPtn = regexp.MustCompile("manifest-.*\\.txt$")
var tagmanifestPtn = regexp.MustCompile("tagmanifest-.*\\.txt$")

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
