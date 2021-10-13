package go_bagit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecksums(t *testing.T) {
	t.Run("Test validate sha256", func(t *testing.T) {
		testFile, err := os.Open(filepath.Join("test", "valid", "data", "test-file.txt"))
		if err != nil {
			t.Error(err)
		}
		checksum := "20cd2eb771177035f483363951203be7cd85f176aaa7d124a56eb4c83562a861"
		if err := ValidateChecksum(testFile, "sha256", checksum); err != nil {
			t.Error(err)
		}
	})
}
