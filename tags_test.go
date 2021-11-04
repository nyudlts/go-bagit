package go_bagit

import (
	"testing"
)

func TestTags(t *testing.T) {
	bagit := CreateBagit()
	t.Run("Test Adding Standard Tag", func(t *testing.T) {
		bagit.Tags[StandardTags.ContactName] = "Donald"
		want := "Donald"
		got := bagit.Tags[StandardTags.ContactName]
		if want != got {
			t.Errorf("Wanted %s got %s", want, got)
		}
	})
}
