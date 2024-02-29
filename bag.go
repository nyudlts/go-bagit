package go_bagit

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Bag struct {
	Path      string
	Name      string
	AbsPath   string
	Payload   interface{}
	TagSets   []TagSet
	Manifests []Manifest
}

func GetExistingBag(path string) (Bag, error) {
	bag := Bag{}

	//check directory exists
	if err := directoryExists(path); err != nil {
		return bag, err
	}

	//set path
	bag.Path = path

	//set absolute path
	var err error
	bag.AbsPath, err = getABS(bag.Path)
	if err != nil {
		return bag, err
	}

	//set name
	pathSplit := strings.Split(bag.AbsPath, string(os.PathSeparator))
	bag.Name = pathSplit[len(pathSplit)-1]

	if err := bag.ValidateBag(false, false); err != nil {
		return bag, err
	}

	return bag, nil
}

func (b Bag) GetAbsolutePath() (string, error) {
	abs, err := filepath.Abs(b.Path)
	if err != nil {
		return abs, err
	}
	return abs, nil
}

func (b Bag) String() string {
	pathSplit := strings.Split(b.AbsPath, string(os.PathSeparator))
	bagPath := strings.Join(pathSplit[:len(pathSplit)-1], string(os.PathSeparator))
	return fmt.Sprintf("%s: %s\n", b.Name, bagPath)
}

type getFilesOrDirsParams struct {
	Location    string
	Matcher     *regexp.Regexp
	FindFiles   bool
	ReturnFirst bool
}

func (b Bag) ValidateBag(fast bool, complete bool) error {
	errs := []error{}
	storedOxum, err := GetOxum(b.Path)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}

	err = ValidateOxum(b.Path, storedOxum)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}

	if fast {
		log.Printf("- INFO - %s valid according to Payload Oxum", b.Path)
		return nil
	}

	//validate ant manifest files
	bagFiles, err := os.ReadDir(b.Path)
	if err != nil {
		return err
	}

	dataFiles := map[string]bool{}
	for _, bagFile := range bagFiles {
		if tagmanifestPtn.MatchString(bagFile.Name()) {
			manifestLoc := filepath.Join(b.Path, bagFile.Name())
			_, e := ValidateManifest(manifestLoc, complete)
			if len(e) > 0 {
				errs = append(errs, e...)
				errorMsgs := gatherErrors(errs, b.Path)
				return errors.New(errorMsgs)
			}
		}

		if manifestPtn.MatchString(bagFile.Name()) {
			manifestLoc := filepath.Join(b.Path, bagFile.Name())
			entries, e := ValidateManifest(manifestLoc, complete)
			if len(e) > 0 {
				errs = append(errs, e...)
			}
			for path := range entries {
				dataFiles[path] = true
			}
		}

	}

	dataDirName := filepath.Join(b.Path, "data")
	if err := filepath.WalkDir(dataDirName, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || dataDirName == path {
			return nil
		}
		rel, err := filepath.Rel(b.Path, path)
		if err != nil {
			return err
		}
		if _, ok := dataFiles[rel]; !ok {
			return fmt.Errorf("%s exists on filesystem but is not in the manifest", rel)
		}
		return nil
	}); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		log.Printf("- INFO - %s is valid", b.Name)
		return nil
	}

	errorMsgs := gatherErrors(errs, b.Name)
	return errors.New(errorMsgs)
}

func gatherErrors(errs []error, bagLocation string) string {
	errorMsgs := fmt.Sprintf("- ERROR - %s is invalid: Bag validation failed: ", bagLocation)
	for i, e := range errs {
		errorMsgs = errorMsgs + e.Error()
		if i < len(errs)-1 {
			errorMsgs = errorMsgs + "; "
		}
	}
	log.Println(errorMsgs)
	return errorMsgs
}

func CreateBag(inputDir string, algorithm string, numProcesses int) error {
	//check that input exists and is a directory
	if err := directoryExists(inputDir); err != nil {
		return err
	}

	log.Printf("- INFO - Creating Bag for directory %s", inputDir)

	//create a slice of files
	filesToBag, err := os.ReadDir(inputDir)
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
	bagit := CreateBagit()
	bagit.Path = inputDir
	/*
		if err := bagit.Serialize(); err != nil {
			return err
		}
	*/

	//Generate bag-info.txt
	log.Println("- INFO - Creating bag-info.txt")

	//get the oxum
	oxum, err := CalculateOxum(inputDir)
	if err != nil {
		return err
	}
	bagInfo := CreateBagInfo()
	bagInfo.Tags[StandardTags.PayloadOxum] = oxum.String()
	bagInfo.Path = inputDir
	/*
		if err := bagInfo.Serialize(); err != nil {
			return err
		}
	*/
	//Generate TagManifest
	if err := CreateTagManifest(inputDir, algorithm, numProcesses); err != nil {
		return err
	}

	//you are done
	return nil
}

// Adds a file to the bag root and registers it in the tag manifest file
func (b Bag) AddFileToBagRoot(file string) error {

	//check if source file is valid
	if err := fileExists(file); err != nil {
		return err
	}

	//check if there is already a source file with the same name in the bag
	sourceFileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	targetFilePath := filepath.Join(b.Path, sourceFileInfo.Name())
	log.Println(targetFilePath)
	err = fileExists(targetFilePath)
	if err == nil {
		return fmt.Errorf("- ERROR - cannot create target file %s already exists", targetFilePath)
	}

	//create the target file
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	//read the source file
	sourceFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	//write the contents of the source file to the target file
	log.Printf("- INFO - copying file %s to %s", file, b.Path)
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}

	//ensure the new file exists
	if err := fileExists(targetFilePath); err != nil {
		return err
	}
	targetFile.Close()

	//locate the tagmanifest
	tagmanifest, err := FindFileInBag(b.Path, regexp.MustCompile("tagmanifest"))
	if err != nil {
		return err
	}

	//append new file to tagmanifest
	if err := appendToTagManifest(targetFilePath, b.Path, filepath.Base(tagmanifest)); err != nil {
		return err
	}

	return nil
}

func GetFilesInBag(bagLocation string) ([]string, error) {
	return getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, regexp.MustCompile(`.*`), true, false})
}

func FindFileInBag(bagLocation string, matcher *regexp.Regexp) (string, error) {
	results, err := getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, matcher, true, true})
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("could not locate file pattern in bag")
	}
	return results[0], nil
}

func FindFilesInBag(bagLocation string, matcher *regexp.Regexp) ([]string, error) {
	return getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, matcher, true, false})
}

func GetDirsInBag(bagLocation string) ([]string, error) {
	return getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, regexp.MustCompile(`.*`), false, false})
}

func FindDirInBag(bagLocation string, matcher *regexp.Regexp) (string, error) {
	results, err := getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, matcher, false, true})
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("could not locate directory pattern in bag")
	}
	return results[0], nil
}

func FindDirsInBag(bagLocation string, matcher *regexp.Regexp) ([]string, error) {
	return getFilesOrDirsInBag(getFilesOrDirsParams{bagLocation, matcher, false, false})
}

// getFilesOrDirsInBag returns a slice of strings of matching files or directories.
// What is returned is controlled by the findFiles boolean.
// findFiles = true   --> return matching files
// findFiles = false  --> return matching directories
//
// How many matches are returned is determined by the returnFirst boolean.
// returnFirst = true  --> halts search and returns with first match
// returnFirst = false --> returns all matching files or directories
func getFilesOrDirsInBag(params getFilesOrDirsParams) ([]string, error) {
	results := []string{}

	bagLocation := params.Location
	matcher := params.Matcher
	findFiles := params.FindFiles
	returnFirst := params.ReturnFirst

	err := filepath.Walk(bagLocation,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// if looking for files, but this is a directory, move on...
			if findFiles && info.IsDir() {
				return nil
			}

			// if looking for directories, but this is NOT a directory, move on...
			if !findFiles && !info.IsDir() {
				return nil
			}

			// OK, we found something that we might be looking for...
			if matcher.MatchString(path) {
				results = append(results, path)
				if returnFirst {
					return filepath.SkipAll
				}
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return results, nil
}
