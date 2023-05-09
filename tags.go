package go_bagit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

func CreateBagInfo() TagSet {
	return TagSet{
		Filename: "bag-info.txt",
		Tags: map[string]string{
			StandardTags.BagSoftwareAgent: GetSoftwareAgent(),
			StandardTags.BaggingDate:      fmt.Sprintf(currentTime.Format("2006-02-01")),
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
		kv := strings.Split(scanner.Text(), ": ")
		tags[kv[0]] = kv[1]
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
}

func GetStandardTags() StandardTagSet {
	standardTags := StandardTagSet{}
	standardTags.BagSoftwareAgent = "Bag_Software_Agent"
	standardTags.SourceOrganization = "Source_Organization"
	standardTags.OrganizationAddress = "Organization_Address"
	standardTags.ContactName = "Contact_Name"
	standardTags.ContactPhone = "Contact_Phone"
	standardTags.ContactEmail = "Contact_email"
	standardTags.ExternalDescription = "External_Description"
	standardTags.ExternalIdentifier = "External_Identifier"
	standardTags.BagSize = "Bag_Size"
	standardTags.BaggingDate = "Bagging_Date"
	standardTags.PayloadOxum = "Payload_Oxum"
	standardTags.BagCount = "Bag_Count"
	standardTags.BagGroupIdentifier = "Bag_Group_Identifier"
	standardTags.InternalSenderIdentifier = "Internal_Sender_Identifier"
	standardTags.InternalSenderDescription = "Internal_Sender_Description"
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
