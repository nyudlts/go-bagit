package go_bagit

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestManifests(t *testing.T) {
	t.Run("Test Parsing Manifest", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		manifestMap, err := ReadManifest(manifestLoc)
		if err != nil {
			t.Error(err)
		}

		want := 1
		got := len(manifestMap)
		if want != got {
			t.Errorf("Wanted %d Got %d", want, got)
		}
	})

	t.Run("Test Validate a manifest file", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		entries, err := ValidateManifest(manifestLoc, false)
		if len(err) > 0 {
			t.Error(err[0])
		}
		want := map[string]string{
			"data/test-file.txt": "20cd2eb771177035f483363951203be7cd85f176aaa7d124a56eb4c83562a861",
		}
		if !reflect.DeepEqual(entries, want) {
			t.Errorf("returned unexpected entries; want: %v, got: %v", want, entries)
		}
	})

	t.Run("Test Completeness Only Validation of a manifest file", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		entries, err := ValidateManifest(manifestLoc, true)
		if len(err) > 0 {
			t.Error(err[0])
		}
		want := map[string]string{
			"data/test-file.txt": "20cd2eb771177035f483363951203be7cd85f176aaa7d124a56eb4c83562a861",
		}
		if !reflect.DeepEqual(entries, want) {
			t.Errorf("returned unexpected entries; want: %v, got: %v", want, entries)
		}
	})
}
