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

func ValidateManifest(manifestLocation string) []error {

	errors := []error{}
	lastInd := strings.LastIndex(manifestLocation, "/")
	path := manifestLocation[:lastInd]

	manifestMap, err := ReadManifest(manifestLocation)
	if err != nil {
		return append(errors, err)
	}

	for k, v := range manifestMap {
		entryPath := filepath.Join(path, k)
		log.Println("- INFO - Verifying checksum for file", entryPath)
		if err := entryExists(entryPath); err != nil {
			return append(errors, err)
		}
		f, err := os.Open(entryPath)
		if err != nil {
			return append(errors, err)
		}

		if err := ValidateSHA256(f, v); err != nil {
			fLocation := f.Name()[len(path)+1 : len(f.Name())]
			err = fmt.Errorf("%s %s", fLocation, err.Error())
			log.Println(fmt.Errorf("- WARNING - %s", err))
			errors = append(errors, err)
		}
	}
	return errors
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
