package go_bagit

import (
	"path/filepath"
	"testing"
)

func TestValidateBag(t *testing.T) {
	t.Run("Test validate a bag", func(t *testing.T) {
		bagLoc := filepath.Join("test", "valid")
		if err := ValidateBag(bagLoc, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test fast validate a bag", func(t *testing.T) {
		bagLoc := filepath.Join("test", "valid")
		if err := ValidateBag(bagLoc, true); err != nil {
			t.Error(err)
		}
	})
}
