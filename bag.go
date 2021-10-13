package go_bagit

import (
	"fmt"
	"log"
	"path/filepath"
)

func ValidateBag(bagLocation string, fast bool) error {
	errors := []error{}
	storedOxum, err := GetOxum(bagLocation)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}

	err = ValidateOxum(bagLocation, storedOxum)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}
	if fast == true {
		log.Printf("- INFO - %s valid according to Payload Oxum", bagLocation)
		return nil
	}

	manifest := filepath.Join(bagLocation, "manifest-sha256.txt")
	e := ValidateManifest(manifest)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	tagmanifest := filepath.Join(bagLocation, "tagmanifest-sha256.txt")
	e = ValidateManifest(tagmanifest)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	if len(errors) == 0 {
		log.Printf("- INFO - %s is valid", bagLocation)
	} else {
		errorMsgs := fmt.Sprintf("- ERROR - %s is invalid: Bag validation failed: ", bagLocation)
		for i, e := range errors {
			errorMsgs = errorMsgs + e.Error()
			if i < len(errors)-1 {
				errorMsgs = errorMsgs + "; "
			}
		}
		log.Println(errorMsgs)
	}
	return nil
}
