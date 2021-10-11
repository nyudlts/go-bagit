package go_bagit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Oxum struct {
	Size int64
	Count int
}

func (o Oxum) String() string {
	return fmt.Sprintf("%d.%d", o.Size, o.Count)
}

func ParseOxumString(oxum string) (Oxum, error) {
	o := Oxum{}
	splitOxum  := strings.Split(oxum, ".")
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


func ValidateOxum(bagLocation string, oxum string) error {
	storedOxum, err := ParseOxumString(oxum)
	if err != nil {
		return err
	}

	log.Println(storedOxum.String())

	calculatedOxum, err := CalculateOxum(bagLocation)
	if err != nil {
		return err
	}

	if calculatedOxum.Size != storedOxum.Size {
		return fmt.Errorf("Size mismatch, expected %v got %v", storedOxum.Size, calculatedOxum.Size)
	}

	if calculatedOxum.Count != storedOxum.Count {
		return fmt.Errorf("Count mismatch, expected %d got %d", storedOxum.Count, calculatedOxum.Count)
	}

	return nil
}

func CalculateOxum(bagLocation string) (Oxum, error){
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
