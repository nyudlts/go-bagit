package go_bagit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var manifestPtn = regexp.MustCompile("manifest-.*\\.txt$")

func ValidateBag(bagLocation string, fast bool, complete bool) error {
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

	//validate ant manifest files
	bagFiles, err := ioutil.ReadDir(bagLocation)
	if err != nil {
		return err
	}

	for _, bagFile := range bagFiles {
		if manifestPtn.MatchString(bagFile.Name()) == true {
			manifestLoc := filepath.Join(bagLocation, bagFile.Name())
			e := ValidateManifest(manifestLoc, complete)
			if len(e) > 0 {
				errors = append(errors, e...)
			}
		}
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

func CreateBag(inputDir string, algorithm string, numProcesses int) error {
	//check that input exists and is a directory
	if err := directoryExists(inputDir); err != nil {
		return err
	}

	log.Printf("- INFO - Creating Bag for directory %s", inputDir)

	//create a slice of files
	filesToBag, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	//check there is at least one file to be bagged.
	if len(filesToBag) < 1 {
		errMsg := fmt.Errorf("Could not create a bag, no files present in %s", inputDir)
		log.Println("- ERROR -", errMsg)
		return errMsg
	}

	//create a data directory for payload
	log.Println("- INFO - Creating data directory")
	dataDirName := filepath.Join(inputDir, "data")
	if err := os.Mkdir(dataDirName, 0777); err != nil {
		log.Println("- ERROR -", err)
		return err
	}

	//move the payload files into data dir
	for _, file := range filesToBag {
		originalLocation := filepath.Join(inputDir, file.Name())
		newLocation := filepath.Join(dataDirName, file.Name())
		log.Printf("- INFO - Moving %s to %s", originalLocation, newLocation)
		if err := os.Rename(originalLocation, newLocation); err != nil {
			log.Println("- ERROR -", err.Error())
			return err
		}
	}

	//Generate the manifest
	if err := CreateManifest("manifest", inputDir, algorithm, numProcesses); err != nil {
		return err
	}

	//Generate bagit.txt
	log.Println("- INFO - Creating bagit.txt")
	if err := createBagIt(inputDir); err != nil {
		return err
	}

	//Generate bag-info.txt
	log.Println("- INFO - Creating bag-info.txt")
	if err := createBagInfo(inputDir); err != nil {
		return err
	}

	//Generate TagManifest
	if err := CreateTagManifest(inputDir, algorithm, numProcesses); err != nil {
		return err
	}

	//you are done
	return nil
}

func createBagInfo(bagLoc string) error {

	oxum, err := CalculateOxum(bagLoc)
	if err != nil {
		return err
	}

	bagInfo.Tags["Payload-Oxum"] = oxum.String()
	bagInfo.Path = bagLoc
	if err := bagInfo.Serialize(); err != nil {
		return err
	}
	return nil
}

func createBagIt(bagLoc string) error {
	bagit.Path = bagLoc
	if err := bagit.Serialize(); err != nil {
		return err
	}
	return nil
}

func directoryExists(inputDir string) error {
	if fi, err := os.Stat(inputDir); err == nil {
		if fi.IsDir() == true {
			return nil
		} else {
			errorMsg := fmt.Errorf("Failed to create bag: input directory %s is not a directory", inputDir)
			log.Println("- Error -", errorMsg)
			return errorMsg
		}
	} else if os.IsNotExist(err) {
		errorMsg := fmt.Errorf(" - ERROR - input %s directory does not exist", inputDir)
		log.Println(errorMsg)
		return err
	} else {
		return err
	}
}
