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
	t.Run("Test GetFilesInBag()", func(t *testing.T) {
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
			t.Fatal(err)
		}

		if len(want) != len(got) {
			t.Fatal("length of returned slice does not match expectations")
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

func TestGetDirsInBag(t *testing.T) {
	t.Run("Test GetDirsInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-erecord-with-subdirs")

		want := []string{
			"test/valid-erecord-with-subdirs",
			"test/valid-erecord-with-subdirs/data",
			"test/valid-erecord-with-subdirs/data/logs",
			"test/valid-erecord-with-subdirs/data/logs/transfers",
			"test/valid-erecord-with-subdirs/data/logs/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5",
			"test/valid-erecord-with-subdirs/data/logs/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5/logs",
			"test/valid-erecord-with-subdirs/data/objects",
			"test/valid-erecord-with-subdirs/data/objects/cuid39675",
			"test/valid-erecord-with-subdirs/data/objects/metadata",
			"test/valid-erecord-with-subdirs/data/objects/metadata/transfers",
			"test/valid-erecord-with-subdirs/data/objects/metadata/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5",
			"test/valid-erecord-with-subdirs/data/objects/submissionDocumentation",
			"test/valid-erecord-with-subdirs/data/objects/submissionDocumentation/transfer-fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5",
		}

		got, err := GetDirsInBag(bagRoot)
		if err != nil {
			t.Fatal(err)
		}

		if len(want) != len(got) {
			t.Fatal("length of returned slice does not match expectations")
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
			t.Fatal(msg)
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

func TestFindFilesInBag(t *testing.T) {
	t.Run("Test FindFilesInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-erecord-with-subdirs")

		want := []string{
			"test/valid-erecord-with-subdirs/fales_mss2023_cuid39675_aspace_wo.tsv",
			"test/valid-erecord-with-subdirs/data/objects/metadata/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5/fales_mss2023_cuid39675_aspace_wo.tsv",
		}
		wantPtn := regexp.MustCompile("_aspace_wo.tsv$")

		got, err := FindFilesInBag(bagRoot, wantPtn)
		if err != nil {
			t.Error(err)
		}

		if len(want) != len(got) {
			t.Fatal("length of returned slice does not match expectations")
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
			t.Fatal(msg)
		}
	})

}

func TestFindDirInBag(t *testing.T) {
	t.Run("Test FindDirInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-erecord-with-subdirs")

		want := "test/valid-erecord-with-subdirs/data/objects/cuid39675"
		wantPtn := regexp.MustCompile("objects/cuid39675$")

		got, err := FindDirInBag(bagRoot, wantPtn)
		if err != nil {
			t.Error(err)
		}

		if want != got {
			t.Errorf("\n%v !=\n%v", want, got)
		}
	})
}
