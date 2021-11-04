package go_bagit

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var bagit = TagSet {
	Filename: "bagit.txt",
	Tags: map[string]string {
		"BagIt-Version": "0.97",
		"Tag-File-Character-Encoding": "UTF-8",
	},
}

var bagInfo = TagSet {
	Filename: "bag-info.txt",
	Tags: map[string]string{
		"Bag-Software-Agent": fmt.Sprintf("go-bagit %s <https://github.com/nyudlts/go-bagit>", libraryVersion),
		"Bagging-Date":      fmt.Sprintf(currentTime.Format("2006-02-01")),
	},
}

type TagSet struct {
	Filename string
	Path string
	Tags map[string]string
}

func (t TagSet) Serialize() error {
	outfile := filepath.Join(t.Path, t.Filename)
	tags := []byte{}
	for k,v := range t.Tags {
		tags = append(tags, []byte(fmt.Sprintf("%s: %s\n", k, v))...)
	}
	if err := ioutil.WriteFile(outfile, tags, 0777); err != nil {
		return err
	}
	return nil
}

/*
var standardTags []Tag {
	BagSoftwareAgent          string `yaml:"Bag_Software_Agent,omitempty"`
	SourceOrganization        string `yaml:"Source_Organization,omitempty"`
	OrganizationAddress       string `yaml:"Organization_Address,omitempty"`
	ContactName               string `yaml:"Contact_Name,omitempty"`
	ContactPhone              string `yaml:"Contact_Phone,omitempty"`
	ContactEmail              string `yaml:"Contact_email,omitempty"`
	ExternalDescription       string `yaml:"External_Description,omitempty"`
	ExternalIdentifier        string `yaml:"External_Identifier,omitempty"`
	BagSize                   string `yaml:"Bag_Size,omitempty"`
	BaggingDate               string `yaml:"Bagging_Date,omitempty"`
	PayloadOxum               string `yaml:"Payload_Oxum,omitempty"`
	BagCount                  string `yaml:"Bag_Count,omitempty"`
	BagGroupIdentifier        string `yaml:"Bag_Group_Identifier,omitempty"`
	InternalSenderIdentifier  string `yaml:"Internal_Sender_Identifier,omitempty"`
	InternalSenderDescription string `yaml:"Internal_Sender_Description,omitempty"`
}
 */

/*
type Bag struct {
	BagItVersion             string	`yaml:"BagIt_Version,omitempty"`
	TagFileCharacterEncoding string `yaml:"Tag_File_Character_Encoding,omitempty"`
}

type Tags struct {
	BagSoftwareAgent          string `yaml:"Bag_Software_Agent,omitempty"`
	SourceOrganization        string `yaml:"Source_Organization,omitempty"`
	OrganizationAddress       string `yaml:"Organization_Address,omitempty"`
	ContactName               string `yaml:"Contact_Name,omitempty"`
	ContactPhone              string `yaml:"Contact_Phone,omitempty"`
	ContactEmail              string `yaml:"Contact_email,omitempty"`
	ExternalDescription       string `yaml:"External_Description,omitempty"`
	ExternalIdentifier        string `yaml:"External_Identifier,omitempty"`
	BagSize                   string `yaml:"Bag_Size,omitempty"`
	BaggingDate               string `yaml:"Bagging_Date,omitempty"`
	PayloadOxum               string `yaml:"Payload_Oxum,omitempty"`
	BagCount                  string `yaml:"Bag_Count,omitempty"`
	BagGroupIdentifier        string `yaml:"Bag_Group_Identifier,omitempty"`
	InternalSenderIdentifier  string `yaml:"Internal_Sender_Identifier,omitempty"`
	InternalSenderDescription string `yaml:"Internal_Sender_Description,omitempty"`
}

func (t Tags) ToYaml() ([]byte, error) {
	yamlTags,err := yaml.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	return yamlTags, nil
}

func (b Bag) ToYaml() ([]byte, error) {
	yamlBag, err := yaml.Marshal(b)
	if err != nil {
		return []byte{}, err
	}
	return yamlBag, nil
}
*/


