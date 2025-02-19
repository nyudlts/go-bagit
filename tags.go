package go_bagit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Bagit TagSet
type BagInfo TagSet

var StandardTags = GetStandardTags()

func CreateBagit() TagSet {
	return TagSet{
		Filename: "bagit.txt",
		Tags: map[string]string{
			"BagIt-Version":               "0.97",
			"Tag-File-Character-Encoding": "UTF-8",
		},
	}
}

func CreateBagInfo(t time.Time) TagSet {
	return TagSet{
		Filename: "bag-info.txt",
		Tags: map[string]string{
			StandardTags.BagSoftwareAgent: GetSoftwareAgent(),
			StandardTags.BaggingDate:      t.Round(0).Format(time.DateOnly),
		},
	}
}

type TagSet struct {
	Filename string
	Path     string
	Tags     map[string]string
}

func NewTagSet(filename string, bagLocation string) (TagSet, error) {
	tagFileLocation := filepath.Join(bagLocation, filename)
	tags := map[string]string{}
	tagFile, err := os.Open(tagFileLocation)
	if err != nil {
		return TagSet{}, err
	}

	scanner := bufio.NewScanner(tagFile)
	for scanner.Scan() {
		// Skip blank lines
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		before, after, _ := strings.Cut(scanner.Text(), ": ")
		tags[before] = after
	}

	tagFile.Close()
	return TagSet{filename, bagLocation, tags}, nil
}

func (tagSet TagSet) GetTagSetAsByteSlice() []byte {
	keys := make([]string, 0)
	for k, _ := range tagSet.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tags := []byte{}
	for _, k := range keys {
		tags = append(tags, []byte(fmt.Sprintf("%s: %s\n", k, tagSet.Tags[k]))...)
	}
	return tags
}

type StandardTagSet struct {
	BagSoftwareAgent          string
	SourceOrganization        string
	OrganizationAddress       string
	ContactName               string
	ContactPhone              string
	ContactEmail              string
	ExternalDescription       string
	ExternalIdentifier        string
	BagSize                   string
	BaggingDate               string
	PayloadOxum               string
	BagCount                  string
	BagGroupIdentifier        string
	InternalSenderIdentifier  string
	InternalSenderDescription string
	BagItVersion              string
	TagFileCharacterEncoding  string
}

func GetStandardTags() StandardTagSet {
	standardTags := StandardTagSet{}
	standardTags.BagSoftwareAgent = "Bag-Software-Agent"
	standardTags.SourceOrganization = "Source-Organization"
	standardTags.OrganizationAddress = "Organization-Address"
	standardTags.ContactName = "Contact-Name"
	standardTags.ContactPhone = "Contact-Phone"
	standardTags.ContactEmail = "Contact-email"
	standardTags.ExternalDescription = "External-Description"
	standardTags.ExternalIdentifier = "External-Identifier"
	standardTags.BagSize = "Bag-Size"
	standardTags.BaggingDate = "Bagging-Date"
	standardTags.PayloadOxum = "Payload-Oxum"
	standardTags.BagCount = "Bag-Count"
	standardTags.BagGroupIdentifier = "Bag-Group-Identifier"
	standardTags.InternalSenderIdentifier = "Internal-Sender-Identifier"
	standardTags.InternalSenderDescription = "Internal-Sender-Description"
	standardTags.BagItVersion = "BagIt-Version"
	standardTags.TagFileCharacterEncoding = "Tag-File-Character-Encoding"
	return standardTags
}

func (tagset TagSet) AddTags(newTags map[string]string) {
	for k, v := range newTags {
		tagset.Tags[k] = v
	}
}

func (tagSet TagSet) HasTag(key string) bool {
	for k, _ := range tagSet.Tags {
		if k == key {
			return true
		}
	}
	return false
}

func (tagSet TagSet) UpdateTagFile(key string, value string) {
	for k, _ := range tagSet.Tags {
		if k == key {
			tagSet.Tags[k] = value
		}
	}
}
