package go_bagit

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBagStruct(t *testing.T) {

	t.Run("It will error on invalid directory", func(t *testing.T) {
		_, err := NewBag("test/invalid-bag-path")
		if err != nil {
			t.Log(err)
		} else {
			t.Errorf("Created a bag for an invalid path")
		}
	})

	t.Run("It can create a Bag type", func(t *testing.T) {
		bag, err := NewBag("test/valid")
		if err != nil {
			t.Errorf(err.Error())
		}
		want := 2
		got := len(bag.Manifests)
		if want != got {
			t.Errorf("wanted %d got %d", want, got)
		}
	})
}

func TestValidateBag(t *testing.T) {
	tests := map[string]struct {
		loc  string
		fast bool
		err  string
	}{
		"It validates a bag": {
			loc:  "test/valid",
			fast: false,
		},
		"It validates a bag in fast mode": {
			loc:  "test/valid",
			fast: true,
		},
		"It validates a bag with subdirs": {
			loc:  "test/valid-with-subdirs",
			fast: false,
		},
		"It validates a bag with subdirs in fast mode": {
			loc:  "test/valid-with-subdirs",
			fast: true,
		},
		"It identifies an invalid bag with unexpected files": {
			loc:  "test/unexpected-files",
			fast: false,
			err:  "Bag validation failed: data/test-file.txt exists on filesystem but is not in the manifest",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path, _ := filepath.Abs(tc.loc)
			err := ValidateBag(path, tc.fast, false)

			if tc.err == "" && err != nil {
				t.Fatalf("expected to pass; got: %v", err)
			}

			if tc.err != "" {
				if err == nil {
					t.Fatalf("expected to fail (%s); got nil err", err)
				}
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("error mismatch; want %s, got %v", tc.err, err)
				}
			}
		})
	}
}
