package go_bagit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Payload map[string]os.FileInfo

type PayloadMatch struct {
	Path     string
	FileInfo os.FileInfo
}

func loadPayload(bag *Bag) error {
	dataDir := filepath.Join(bag.Path, "data")
	err := filepath.Walk(dataDir, func(path string, info fs.FileInfo, err error) error {
		if path != dataDir {
			bag.Payload[path] = info
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p Payload) FindFileInPayload(matcher *regexp.Regexp) (PayloadMatch, error) {
	for path, fi := range p {
		if matcher.MatchString(path) && !fi.IsDir() {
			return PayloadMatch{path, fi}, nil
		}
	}
	return PayloadMatch{}, fmt.Errorf("Payload did not match %s", matcher.String())
}
