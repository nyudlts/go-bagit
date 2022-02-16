package go_bagit

import (
	"fmt"
	"time"
)

var currentTime = time.Now()

var libraryVersion = "v0.1.1-alpha"

func GetSoftwareAgent() string {
	return fmt.Sprintf("go-bagit %s <https://github.com/nyudlts/go-bagit>", libraryVersion)
}
