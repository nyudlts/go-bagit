package go_bagit

import (
	"path/filepath"
	"testing"
)

func TestValidateBag(t *testing.T) {
	t.Run("Test validate a bag", func(t *testing.T) {
		bagLoc := filepath.Join("test", "valid")
		if err := ValidateBag(bagLoc, false, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test fast validate a bag", func(t *testing.T) {
		bagLoc := filepath.Join("test", "valid")
		if err := ValidateBag(bagLoc, true, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test validate a bag with unexpected files fails", func(t *testing.T) {
		bagLoc := filepath.Join("test", "unexpected-files")
		if err := ValidateBag(bagLoc, false, false); err == nil {
			t.Error("expected to fail but returned a nil error")
		}
	})
}
