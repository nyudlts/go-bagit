package go_bagit

import (
	"fmt"
	"os"
	"path/filepath"
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

func getABS(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return abs, err
	}
	return abs, nil
}

func fileExists(file string) error {
	if _, err := os.Stat(file); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		errorMsg := fmt.Errorf("file %s does not exist", file)
		return errorMsg
	} else {
		log.Println("- ERROR - unknown error:", err.Error())
		return err
	}
}

func directoryExists(inputDir string) error {
	if fi, err := os.Stat(inputDir); err == nil {
		if fi.IsDir() {
			return nil
		} else {
			errorMsg := fmt.Errorf("- ERROR - input directory %s is not a directory", inputDir)
			return errorMsg
		}
	} else if os.IsNotExist(err) {
		errorMsg := fmt.Errorf("- ERROR - input %s directory does not exist", inputDir)
		return errorMsg
	} else {
		return err
	}
}
