package go_bagit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	cp "github.com/otiai10/copy"
	"gotest.tools/v3/assert"
	tfs "gotest.tools/v3/fs"
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
			bag, _ := GetExistingBag(tc.loc)
			err := bag.ValidateBag(tc.fast, false)

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

func TestFindDirsInBag(t *testing.T) {
	t.Run("Test FindDirsInBag()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-erecord-with-subdirs")

		want := []string{
			"test/valid-erecord-with-subdirs/data/logs/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5",
			"test/valid-erecord-with-subdirs/data/objects/metadata/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5",
		}
		wantPtn := regexp.MustCompile("/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5$")

		got, err := FindDirsInBag(bagRoot, wantPtn)
		if err != nil {
			t.Error(err)
		}

		if len(want) != len(got) {
			t.Fatalf("length of returned slice, (%d) does not match expectations (%d)", len(got), len(want))
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

func TestAddFileToBagRoot(t *testing.T) {
	t.Run("Add a file to bag root", func(t *testing.T) {
		source := filepath.Join("test", "valid")
		target := filepath.Join("test", "addFile", "valid")
		if err := cp.Copy(source, target); err != nil {
			t.Error(err)
		}

		addFile := filepath.Join("test", "addFile", "addfile.txt")
		bag, err := GetExistingBag(target)
		if err != nil {
			t.Error(err)
		}

		if err := bag.AddFileToBagRoot(addFile); err != nil {
			t.Error(err)
		}

		if err := os.RemoveAll(target); err != nil {
			t.Error(err)
		}
	})
}

func TestManifestFunctionsInBag(t *testing.T) {
	t.Run("Test Bag has Manifest File", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(bag.Manifests)
		if len(bag.Manifests) != 1 {
			t.Error("Bag did not contain a manifest file")
		}
	})

	t.Run("Test Bag has Tag Manifest File", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid"))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(bag.TagManifests)
		if len(bag.TagManifests) != 1 {
			t.Error("Bag did not contain a manifest file")
		}
	})
}

func TestTagFileFunctionsInBag(t *testing.T) {
	t.Run("Test Get Bagit.txt file", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid"))
		if err != nil {
			t.Fatal(err)
		}

		tags := map[string]string{
			StandardTags.BagItVersion:             "0.97",
			StandardTags.TagFileCharacterEncoding: "UTF-8",
		}

		for k, v := range tags {
			got := bag.BagIt.Tags[k]
			want := v
			if want != got {
				t.Errorf("Wanted %s got %s", want, got)
			}
		}
	})

	t.Run("Test Get Bag-info.txt file", func(t *testing.T) {
		bag, err := GetExistingBag(filepath.Join("test", "valid"))
		if err != nil {
			t.Fatal(err)
		}

		tags := map[string]string{
			StandardTags.BagSoftwareAgent: "bagit.py v1.8.1 <https://github.com/LibraryOfCongress/bagit-python>",
			StandardTags.BaggingDate:      "2021-10-11",
		}

		for k, v := range tags {
			got := bag.BagInfo.Tags[k]
			want := v
			if want != got {
				t.Errorf("Wanted %s got %s", want, got)
			}
		}
	})
}

func TestCreateABag(t *testing.T) {
	var (
		baginfoContent = fmt.Sprintf(
			`Bag-Software-Agent: go-bagit (%s) <https://github.com/nyudlts/go-bagit>
Bagging-Date: %s
Payload-Oxum: 42.2
`,
			version,
			time.Now().Round(0).Format(time.DateOnly),
		)

		bagitContent = `BagIt-Version: 0.97
Tag-File-Character-Encoding: UTF-8
`
		manifestContent = `bf73d81371ea21348bfb510d8c8948bb64e0eb3cea97ec991a4170e777b6de18  data/test.txt
091e8c6b2cb96659e3c52b3b585d0885989bba9eed34d77563fe0abaf64741ea  data/test2.txt
`

		tagManifestContent = `0391a7a6394036667313e14571bf26e10b8c1d4329b82cb919a12b7936c549ae  bag-info.txt
e91f941be5973ff71f1dccbdd1a32d598881893a7f21be516aca743da38b1689  bagit.txt
ea937667fac06dad61c25afa53a134e97b0b51f6b912b8548e4b32ef2dc2b7f6  manifest-sha256.txt
`
	)

	t.Run("Create A Bag From Existing Dir", func(t *testing.T) {
		sourceDir := filepath.Join("test", "bag-me")
		testDir := tfs.NewDir(t, "gobagit_create_bag_test")
		if err := cp.Copy(sourceDir, testDir.Path()); err != nil {
			t.Fatalf("failed to copy test directory: %v", err)
		}

		_, err := CreateBag(testDir.Path(), "sha256", 1)
		if err != nil {
			t.Error(err)
		}

		assert.Assert(t, tfs.Equal(testDir.Path(), tfs.Expected(t,
			tfs.WithMode(0o755),
			tfs.WithFile("bag-info.txt", baginfoContent, tfs.WithMode(fileMode)),
			tfs.WithFile("bagit.txt", bagitContent, tfs.WithMode(fileMode)),
			tfs.WithFile("manifest-sha256.txt", manifestContent, tfs.WithMode(fileMode)),
			tfs.WithFile("tagmanifest-sha256.txt", tagManifestContent, tfs.WithMode(fileMode)),
			tfs.WithDir("data", tfs.WithMode(dirMode),
				tfs.WithFile("test.txt", "I am a test file.\n", tfs.WithMode(fileMode)),
				tfs.WithFile("test2.txt", "I am another test file.\n", tfs.WithMode(fileMode)),
			),
		)))
	})
}
