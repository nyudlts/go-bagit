package go_bagit

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
			StandardTags.BagSoftwareAgent: fmt.Sprintf("go-bagit %s <https://github.com/nyudlts/go-bagit>", libraryVersion),
			StandardTags.BaggingDate: fmt.Sprintf(currentTime.Format("2006-02-01")),
		},
	}
}

type TagSet struct {
	Filename string
	Path     string
	Tags     map[string]string
}

func (t TagSet) Serialize() error {
	outfile := filepath.Join(t.Path, t.Filename)
	tags := []byte{}
	for k, v := range t.Tags {
		tags = append(tags, []byte(fmt.Sprintf("%s: %s\n", k, v))...)
	}
	if err := ioutil.WriteFile(outfile, tags, 0777); err != nil {
		return err
	}
	return nil
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
