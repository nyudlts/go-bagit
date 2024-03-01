package go_bagit

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestPayloadFuncs(t *testing.T) {
	t.Run("Test GetFileInPayload()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-with-subdirs")
		bag, err := GetExistingBag(bagRoot)
		if err != nil {
			t.Error(err)
		}
		want := "output2.log"

		got, err := bag.Payload.GetFileInPayload(want)
		if err != nil {
			t.Error(err)
		}

		if want != got.FileInfo.Name() {
			t.Errorf("\n%v !=\n%v", want, got)
		}
	})
}

func TestFindFilesInPayload(t *testing.T) {
	t.Run("Test FindFilesInPayload()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-erecord-with-subdirs")
		bag, err := GetExistingBag(bagRoot)
		if err != nil {
			t.Error(err)
		}
		want := 1
		wantPtn := regexp.MustCompile("_aspace_wo.tsv$")

		got := bag.Payload.FindFilesInPayload(wantPtn)
		if err != nil {
			t.Error(err)
		}

		if want != len(got) {
			t.Fatal("length of returned slice does not match expectations")
		}

		if got[0].Path != "test/valid-erecord-with-subdirs/data/objects/metadata/transfers/fales_mss2023_cuid39675-48b63462-0fec-4f6a-8913-1f2e2f9168e5/fales_mss2023_cuid39675_aspace_wo.tsv" {
			t.Errorf("Path value does not match expectations")
		}

	})

}
