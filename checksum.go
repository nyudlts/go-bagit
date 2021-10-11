package go_bagit

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func ValidateSHA256(f *os.File, checksum string) error {
	calculatedChecksum, err := GenerateSha256Checksum(f)
	if err != nil {
		return err
	}

	if calculatedChecksum != checksum {
		return fmt.Errorf("Checksum did not math")
	}
	return nil
}

func GenerateSha256Checksum(f *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
