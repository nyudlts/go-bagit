package go_bagit

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ReadManifest(path string) (map[string]string, error) {
	manifestEntryMap := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return manifestEntryMap, err
	}
	r := regexp.MustCompile("[^\\s]+")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		entry := r.FindAllString(line, -1)
		manifestEntryMap[entry[1]] = entry[0]
	}

	return manifestEntryMap, nil

}

func appendToTagManifest(targetFilePath string, bagLocation string, manifestFileName string) error {

	//get target file metadata
	targetFileInfo, err := os.Stat(targetFilePath)
	if err != nil {
		return err
	}
	log.Printf("- INFO - Adding %s to %s", targetFileInfo.Name(), manifestFileName)

	//get Algorithm of manifest file
	manifestLocation := filepath.Join(bagLocation, manifestFileName)
	algorithm := getAlgorithm(manifestLocation)

	//generate the checksum
	log.Printf("- INFO - generating checksum for %s using %s algorithm", targetFileInfo.Name(), algorithm)
	targetFile, err := os.Open(targetFilePath)
	checksum, err := GenerateChecksum(targetFile, algorithm)
	if err != nil {
		return err
	}
	entry := fmt.Sprintf("%s  %s", checksum, targetFileInfo.Name())

	//open the manifest file
	manifestFile, err := os.OpenFile(manifestLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer manifestFile.Close()

	//write the entry to the manifest file
	if _, err := manifestFile.WriteString(entry); err != nil {
		return err
	}

	return nil
}

func ValidateManifest(manifestLocation string, complete bool) []error {

	errors := []error{}
	lastInd := strings.LastIndex(manifestLocation, "/")
	path := manifestLocation[:lastInd]
	file := manifestLocation[lastInd:]
	algorithm := getAlgorithm(file)
	manifestMap, err := ReadManifest(manifestLocation)
	if err != nil {
		return append(errors, err)
	}

	for k, v := range manifestMap {
		entryPath := filepath.Join(path, k)
		absolutePath, _ := filepath.Abs(entryPath)

		if err := entryExists(entryPath); err != nil {
			return append(errors, err)
		}
		f, err := os.Open(entryPath)
		if err != nil {
			return append(errors, err)
		}

		if complete == false {
			log.Println("- INFO - Verifying checksum for file", absolutePath)
			if err := ValidateChecksum(f, algorithm, v); err != nil {
				fLocation := f.Name()[len(path)+1 : len(f.Name())]
				err = fmt.Errorf("%s %s", fLocation, err.Error())
				log.Println(fmt.Errorf("- WARNING - %s", err))
				errors = append(errors, err)
			}
		}
	}
	return errors
}



func entryExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	} else {
		return err
	}
}

func getAlgorithm(filename string) string {
	split := strings.Split(filename, "-")
	removeExtension := strings.Split(split[len(split) - 1], ".")
	return removeExtension[0]
}

func CreateManifest(manifestName string, bagLoc string, algorithm string, numProcesses int) error {
	dataDir := filepath.Join(bagLoc, "data")
	log.Printf("- INFO - Using %d processes to generate manifests: %s", numProcesses, algorithm)
	manifestLines := []string{}
	err := filepath.WalkDir(dataDir, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() != true {
			log.Printf("- INFO - Generating manifest lines for file %s", path)
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			checksum, err := GenerateChecksum(f, algorithm)
			if err != nil {
				return err
			}
			entryName := path[len(bagLoc)+1:]
			manifestLines = append(manifestLines, fmt.Sprintf("%s  %s", checksum, entryName))
		}
		return nil
	})
	if err != nil {
		return err
	}

	manifestFileName := fmt.Sprintf("%s-%s.txt", manifestName, algorithm)
	if err := createManifestFile(bagLoc, manifestFileName, manifestLines); err != nil {
		return err
	}
	return nil
}

func CreateTagManifest(inputDir string, algorithm string, numProcesses int) error {

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	manifestLines := []string{}

	for _, file := range files {
		if file.IsDir() != true {
			log.Printf("- INFO - Generating manifest lines for file %s", file.Name())
			fi, err := os.Open(filepath.Join(inputDir, file.Name()))
			if err != nil {
				return err
			}
			checksum, err := GenerateChecksum(fi, algorithm)
			if err != nil {
				return err
			}

			manifestLines = append(manifestLines, fmt.Sprintf("%s  %s", checksum, file.Name()))
		}
	}

	manifestName := fmt.Sprintf("tagmanifest-%s.txt", algorithm)

	if err := createManifestFile(inputDir, manifestName, manifestLines); err != nil {
		return err
	}

	return nil
}

func createManifestFile(bagLocation string, manifestFileName string, manifestLines []string) error {
	outFile, err := os.Create(filepath.Join(bagLocation, manifestFileName))
	if err != nil {
		return err
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	for _, manifestLine := range manifestLines {
		writer.WriteString(manifestLine + "\n")
		writer.Flush()
	}
	return nil
}
