package go_bagit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Oxum struct {
	Size  int64
	Count int
}

func (o Oxum) String() string {
	return fmt.Sprintf("%d.%d", o.Size, o.Count)
}

func ParseOxumString(oxum string) (Oxum, error) {
	o := Oxum{}
	splitOxum := strings.Split(oxum, ".")
	var err error
	o.Size, err = strconv.ParseInt(splitOxum[0], 10, 64)
	if err != nil {
		return o, err
	}

	o.Count, err = strconv.Atoi(splitOxum[1])
	if err != nil {
		return o, err
	}
	return o, nil
}

func GetOxum(bagLocation string) (string, error) {
	f, err := os.Open(filepath.Join(bagLocation, "bag-info.txt"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ":")
		if splitLine[0] == "Payload-Oxum" {
			return strings.TrimSpace(splitLine[1]), nil
		}
	}

	return "", fmt.Errorf("%s did not contain a payload-oxum", f.Name())
}

func ValidateOxum(bagLocation string, oxum string) error {
	storedOxum, err := ParseOxumString(oxum)
	if err != nil {
		return err
	}

	calculatedOxum, err := CalculateOxum(bagLocation)
	if err != nil {
		return err
	}

	if calculatedOxum.Size != storedOxum.Size || calculatedOxum.Count != storedOxum.Count {
		return fmt.Errorf("%s is invalid: Payload-Oxum validation failed. Expected %d files and %v bytes but found %d files and %v bytes",
			bagLocation, storedOxum.Count, storedOxum.Size, calculatedOxum.Count, calculatedOxum.Size)
	}

	return nil
}

func CalculateOxum(bagLocation string) (Oxum, error) {
	dataDir := filepath.Join(bagLocation, "data")

	var size int64 = 0
	var count int = 0

	if err := filepath.Walk(dataDir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
			count += 1
		}
		return err
	}); err != nil {
		return Oxum{}, err
	}

	return Oxum{size, count}, nil
}
