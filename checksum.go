package go_bagit

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
)

func ValidateChecksum(f *os.File, algorithm string, checksum string) error {
	calculatedChecksum, err := GenerateChecksum(f, algorithm)
	if err != nil {
		return err
	}

	if calculatedChecksum != checksum {
		return fmt.Errorf("%s validation failed: expected=\"%s\" found=\"%s\"", algorithm, checksum, calculatedChecksum)
	}
	return nil
}

func GenerateChecksum(f *os.File, algorithm string) (string, error) {
	h := getHash(algorithm)
	if h == nil {
		return "", fmt.Errorf("%s is not a supported alogorithm", algorithm)
	}
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getHash(code string) hash.Hash {
	if code == "md5" {
		return md5.New()
	}

	if code == "sha1" {
		return sha1.New()
	}

	if code == "sha256" {
		return sha256.New()
	}

	if code == "sha512" {
		return sha512.New()
	}

	return nil
}
