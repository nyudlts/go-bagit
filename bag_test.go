package go_bagit

import (
	"path/filepath"
	"testing"
)

func TestValidateBag(t *testing.T) {
	t.Run("Test validate a bag", func(t *testing.T) {
		bagLoc := filepath.Join("test", "valid")
		if err := ValidateBag(bagLoc); err != nil {
			t.Error(err)
		}
	})
}