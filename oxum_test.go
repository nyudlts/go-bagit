package go_bagit

import (
	"path/filepath"
	"testing"
)

var testBag = filepath.Join("test", "valid")

func TestOxum(t *testing.T) {
	t.Run("Test Oxum Calculation", func(t *testing.T) {

		oxum, err := CalculateOxum(testBag)
		if err != nil {
			t.Error(err)
		}
		want := "20.1"
		if(oxum.String() != want) {
			t.Error("oxum did not calculate")
		}
	})

	t.Run("Test Oxum Validation", func(t *testing.T) {

		if err := ValidateOxum(testBag, "20.1"); err != nil {
			t.Error(err)
		}
	})
}
