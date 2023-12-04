package go_bagit

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

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

func TestGetFilesInBag(t *testing.T) {
	t.Run("Test FindFilesInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-with-subdirs")

		want := []string{
			"test/valid-with-subdirs/bagit.txt",
			"test/valid-with-subdirs/bag-info.txt",
			"test/valid-with-subdirs/manifest-sha512.txt",
			"test/valid-with-subdirs/tagmanifest-sha512.txt",
			"test/valid-with-subdirs/data/test-file.txt",
			"test/valid-with-subdirs/data/logs/output2.log",
			"test/valid-with-subdirs/data/logs/output1.log",
		}

		got, err := GetFilesInBag(bagRoot)
		if err != nil {
			t.Error(err)
		}

		if len(want) != len(got) {
			t.Error("length of returned slice does not match expectations")
		}

		slices.Sort(want)
		slices.Sort(got)

		status := true
		msg := ""
		for i := 0; i < len(want); i++ {
			if want[i] != got[i] {
				status = false
				msg = msg + "\n" + fmt.Sprintf("%v != %v", want[i], got[i])
			}
		}
		if !status {
			t.Error(msg)
		}
	})
}

func TestFindFileInBag(t *testing.T) {
	t.Run("Test FindFileInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-with-subdirs")

		want := "test/valid-with-subdirs/data/logs/output2.log"
		wantPtn := regexp.MustCompile("output2.log$")

		got, err := FindFileInBag(bagRoot, wantPtn)
		if err != nil {
			t.Error(err)
		}

		if want != got {
			t.Errorf("\n%v !=\n%v", want, got)
		}
	})
}
