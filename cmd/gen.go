package cmd

import (
	cfg "atc/config"
	pkg "atc/packages"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

// RunGen is called on running cf gen
func (opt Opts) RunGen() {
	// check if any templates exist
	if len(cfg.Templates) == 0 {
		pkg.Log.Error("No configured template's exist")
		return
	}
	// index of template config to use
	idx := cfg.Settings.DfltTmplt
	if len(cfg.Templates) == 1 {
		idx = 0
	} else if idx == -1 || opt.All == true {
		// ask user to select desired template
		err := survey.AskOne(&survey.Select{
			Message: "Select template to generate:",
			Options: cfg.ListTmplts(-1),
		}, &idx)
		pkg.PrintError(err, "")
	}
	// create template in current folder
	// leaving path to "" creates file in curr directory
	opt.GenCode(&cfg.Templates[idx], "")
	return
}

// GenCode is to generate the code file in given path
func (opt Opts) GenCode(t *cfg.Template, path string) {
	// read template code file
	file, err := ioutil.ReadFile(t.Path)
	pkg.PrintError(err, "Failed to read template file")
	// clean template code (replace placeholders)
	e := Env{
		Contest: opt.contest,
		Problem: opt.problem,
	}

	source := e.ReplPlaceholder(string(file))

	// name of file to be created
	fName := fmt.Sprintf("${problem}${idx}%v", t.Ext)
	for idx := 0; ; idx++ {
		// idx value to replace in string
		e.Idx = strconv.Itoa(idx)
		name := e.ReplPlaceholder(fName)

		// check if file already exists
		if _, err := os.Stat(name); os.IsNotExist(err) {
			pkg.CreateFile(source, filepath.Join(path, name))
			pkg.Log.Notice("File " + name + " generated")
			break
		}
		pkg.Log.Warning("File " + name + " exists")
	}
	return
}
