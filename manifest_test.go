package go_bagit

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestManifests(t *testing.T) {
	t.Run("Test Parsing Manifest", func(t *testing.T) {
		manifestLoc := filepath.Join("test", "valid", "manifest-sha256.txt")
		manifestMap, err := ReadManifestMap(manifestLoc)
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

	t.Run("Test Get Manifests", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid-with-subdirs"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(bag)

		manifests, err := GetManifests(bag.Path)
		if err != nil {
			t.Error(err)
		}

		t.Log(manifests)
		if len(manifests) != 1 {
			t.Errorf("ERROR")
		}
	})

	t.Run("Test Get Multiple Manifests", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid-with-multi-manifests"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(bag)

		manifests, err := GetManifests(bag.Path)
		if err != nil {
			t.Error(err)
		}

		t.Log(manifests)
		if len(manifests) != 2 {
			t.Errorf("ERROR")
		}
	})

	t.Run("Test Get Tag Manifests", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid-with-subdirs"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(bag)

		manifests, err := GetTagManifests(bag.Path)
		if err != nil {
			t.Error(err)
		}

		t.Log(manifests)
		if len(manifests) != 1 {
			t.Errorf("ERROR")
		}
	})

	t.Run("Test Get Multiple Tag Manifests", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid-with-multi-manifests"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(bag)

		manifests, err := GetTagManifests(bag.Path)
		if err != nil {
			t.Error(err)
		}

		t.Log(manifests)
		if len(manifests) != 2 {
			t.Errorf("ERROR")
		}
	})
}
