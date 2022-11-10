package go_bagit

import "testing"

var testBagLoc = "/var/archivematica/sharedDirectory/www/AIPStore/0a3a/994e/4f60/494a/9f98/5b36/acfa/8cbb/TW_TAM_235-0a3a994e-4f60-494a-9f98-5b36acfa8cbb/"

func TestBag(t *testing.T) {
	var bag *Bag
	var err error
	t.Run("Test getting bag from FS", func(t *testing.T) {
		bag, err = GetBag(testBagLoc)
		if err != nil {
			t.Errorf(err.Error())
		}
	})

	t.Run("Test validating bag", func(t *testing.T) {
		err := bag.ValidateBag()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test printing bag info", func(t *testing.T) {
		t.Log(bag)
	})
}

/*
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
			path, _ := filepath.Abs(tc.loc)
			err := ValidateBag(path, tc.fast, false)

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

*/
