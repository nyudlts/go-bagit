package go_bagit

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestTags(t *testing.T) {
	t.Parallel()

	t.Run("Test Adding Standard Tag", func(t *testing.T) {
		t.Parallel()

		bagit := CreateBagit()
		bagit.Tags[StandardTags.ContactName] = "Donald"
		want := "Donald"
		got := bagit.Tags[StandardTags.ContactName]
		if want != got {
			t.Errorf("Wanted %s got %s", want, got)
		}
	})

	t.Run("Test Querying Non-existant Tag", func(t *testing.T) {
		t.Parallel()

		bagit := CreateBagit()
		want := ""
		got := bagit.Tags["daea5275-bac0-486e-8cac-f1a061c623f6"]
		if want != got {
			t.Errorf("Wanted %s got %s", want, got)
		}
	})
}

func TestCreateBagInfo(t *testing.T) {
	t.Parallel()

	bi := CreateBagInfo(time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC))
	assert.DeepEqual(t, bi, TagSet{
		Filename: "bag-info.txt",
		Tags: map[string]string{
			StandardTags.BagSoftwareAgent: fmt.Sprintf(
				"go-bagit (%s) <https://github.com/nyudlts/go-bagit>", version,
			),
			StandardTags.BaggingDate: "2024-05-15",
		},
	})
}

func TestParseBagInfoWithBlankLines(t *testing.T) {
	t.Parallel()

	ts, err := NewTagSet("bag-info.txt", filepath.Join(".", "test", "valid-bag-info-with-blank-lines"))
	assert.NilError(t, err)
	assert.DeepEqual(t, ts, TagSet{
		Filename: "bag-info.txt",
		Path:     "test/valid-bag-info-with-blank-lines",
		Tags: map[string]string{
			"Bag-Software-Agent": "bagit.py v1.8.1 <https://github.com/LibraryOfCongress/bagit-python>",
			"Bagging-Date":       "2021-10-11",
			"Payload-Oxum":       "20.1",
			"nyudl-comment":      "this is a comment after some blank lines",
		},
	})
}
