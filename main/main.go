package main

import (
	"os"

	go_bagit "github.com/nyudlts/go-bagit"
	"github.com/nyudlts/go-bagit/cmd"
)

func main() {
	go_bagit.Logger().SetOutput(os.Stderr)

	cmd.Execute()
}
