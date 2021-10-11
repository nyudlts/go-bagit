package go_bagit

import (
	"bufio"
	"fmt"
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

	storedOxum, err := GetOxum(bagLocation)
	if err != nil {
		return err
	}

	err = ValidateOxum(bagLocation, storedOxum)
	if err != nil {
		return err
	}

	manifest := filepath.Join(bagLocation, "manifest-sha256.txt")
	if err := ValidateManifest(manifest); err != nil {
		return err
	}
	return nil
}
