package go_bagit

import (
	"regexp"
	"testing"
)

func TestGetSoftwareAgent(t *testing.T) {
	agent := GetSoftwareAgent()
	agentPtn := regexp.MustCompile("go-bagit.*<https://github.com/nyudlts/go-bagit>")
	if agentPtn.MatchString(agent) == false {
		t.Error(agent)
	}
}

func TestCommon(t *testing.T) {
	t.Run("Directory Exists", func(t *testing.T) {
		err := DirectoryExists("test/valid")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Directory Does Not Exists", func(t *testing.T) {
		err := DirectoryExists("test/does-not-exist")
		if err == nil {
			t.Error(err)
		}
	})

	t.Run("Directory Is A File", func(t *testing.T) {
		err := DirectoryExists("test/baginfo.txt")
		if err == nil {
			t.Error(err)
		}
	})
}
