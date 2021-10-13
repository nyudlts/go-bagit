package go_bagit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetOxum(bagLocation string) (string, error) {
	f, err := os.Open(filepath.Join(bagLocation, "bag-info.txt"))
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ":")
		if splitLine[0] == "Payload-Oxum" {
			return strings.TrimSpace(splitLine[1]), nil
		}
	}

	return "", fmt.Errorf("bag-info.txt did not contain a payload-oxum")
}

func ValidateBag(bagLocation string) error {
	errors := []error{}
	storedOxum, err := GetOxum(bagLocation)
	if err != nil {
		errors = append(errors, err)
	}

	err = ValidateOxum(bagLocation, storedOxum)
	if err != nil {
		errors = append(errors, err)
	}

	manifest := filepath.Join(bagLocation, "manifest-sha256.txt")
	e := ValidateManifest(manifest)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	tagmanifest := filepath.Join(bagLocation, "tagmanifest-sha256.txt")
	e = ValidateManifest(tagmanifest)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	if len(errors) == 0 {
		log.Printf("- INFO - %s is valid", bagLocation)
	} else {
		errorMsgs := fmt.Sprintf("- ERROR - %s is invalid: Bag validation failed: ", bagLocation)
		for i, e := range errors {
			errorMsgs = errorMsgs + e.Error()
			if i < len(errors)-1 {
				errorMsgs = errorMsgs + "; "
			}
		}
		log.Println(errorMsgs)
	}
	return nil
}
