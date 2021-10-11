package go_bagit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ReadManifest(path string) (map[string]string, error) {
	manifestEntryMap := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return manifestEntryMap, err
	}
	r := regexp.MustCompile("[^\\s]+")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		entry := r.FindAllString(line, -1)
		manifestEntryMap[entry[1]] = entry[0]
	}

	return manifestEntryMap, nil

}

func ValidateManifest(manifestLocation string) error {

	lastInd := strings.LastIndex(manifestLocation, "/")
	path := manifestLocation[:lastInd]

	manifestMap, err := ReadManifest(manifestLocation)
	if err != nil {
		return err
	}

	for k,v  := range manifestMap {
		entryPath := filepath.Join(path, k)
		if err := entryExists(entryPath); err != nil {
			return err
		}
		f, err := os.Open(entryPath)
		if err != nil {
			return err
		}
		log.Println("- INFO - Verifying checksum for file", entryPath)
		if err := ValidateSHA256(f, v); err != nil {
			return 	err
		}
	}
	return nil
}

func entryExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	} else {
		return err
	}
}