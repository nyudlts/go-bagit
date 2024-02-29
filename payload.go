package go_bagit

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func loadPayload(bag *Bag) error {
	dataDir := filepath.Join(bag.Path, "data")
	err := filepath.Walk(dataDir, func(path string, info fs.FileInfo, err error) error {
		if path != dataDir {
			bag.Payload = append(bag.Payload, info)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b Bag) FindFileinPayload(matcher *regexp.Regexp) (string, error) {
	for _, fi := range b.Payload {
		if matcher.MatchString(fi.Name()) && !fi.IsDir() {
			return fi.Name(), nil
		}
	}
	return "", fmt.Errorf("Payload did not match %s", matcher.String())
}
