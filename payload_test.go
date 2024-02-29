package go_bagit

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestPayloadFuncs(t *testing.T) {
	t.Run("Test FindFileInPayload()", func(t *testing.T) {
		bagRoot := filepath.Join("test", "valid-with-subdirs")
		bag, err := GetExistingBag(bagRoot)
		if err != nil {
			t.Error(err)
		}
		want := "test/valid-with-subdirs/data/logs/output2.log"
		wantPtn := regexp.MustCompile("output2.log$")

		got, err := bag.Payload.FindFileInPayload(wantPtn)
		if err != nil {
			t.Error(err)
		}

		if want != got {
			t.Errorf("\n%v !=\n%v", want, got)
		}
	})
}
