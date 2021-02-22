package gomod

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Copied from `go help mod edit`
type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Module struct {
	Path    string
	Version string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

// GetModule obtains information about the module located at `basePath`
//    Assumes that the `go` executable is on the PATH.
func GetModule(basePath string) (*GoMod, error) {
	cmdGoMod := exec.Command("go", "mod", "edit", "-json")
	cmdGoMod.Dir = basePath

	result, err := cmdGoMod.Output()
	if err != nil {
		return nil, fmt.Errorf("obtaining module data (via '%s') failed:%s", cmdGoMod.String(), err)
	}

	goMod := &GoMod{}
	err = json.Unmarshal(result, goMod)
	if err != nil {
		return nil, fmt.Errorf("parsing Go Module file failed: %s", err)
	}

	return goMod, nil
}
