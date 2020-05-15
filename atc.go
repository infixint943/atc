package main

import (
	cmd "atc/cmd"
	cfg "atc/config"

	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

const manPage = `
Usage:
  atc config
  atc gen    [-A]
  atc open   [<info>...]
  atc fetch  [<info>...]
  atc test   [[-i -e<e> -t<t>] | -C] [-f<f>]
  atc submit [<info>... -f<f>]

Options:
  -A, --all                   force the selection menu to appear
  -f, --file <f>              specify source file to test / submit [default: *.*]
  -i, --ignore-case           omit character-case differences in output
  -e, --ignore-exp <e>        omit float differences <= 1e-<e> [default: 10]
  -t, --time-limit <t>        set time limit (secs) for each test case [default: 2]
  -C, --custom                run interactive session, with input from stdin
  -h, --help                  show this screen
  -v, --version               show cli version
`

func main() {
	args, _ := docopt.ParseArgs(manPage, os.Args[1:], cmd.Version)
	// create ~/atc/ folder
	path, _ := os.UserConfigDir()
	path = filepath.Join(path, "atc")
	os.Mkdir(path, os.ModePerm)
	// initialise default values of atc tool
	// WARNING InitSession() depends on InitSettings()
	cfg.InitTemplates(filepath.Join(path, "templates.json"))
	cfg.InitSettings(filepath.Join(path, "settings.json"))
	cfg.InitSession(filepath.Join(path, "sessions.json"))
	// bind data to struct holding flags
	// and extract contest type / path
	opt := cmd.Opts{}
	args.Bind(&opt)
	opt.FindContestData()
	// pkg.IsAPI(opt.API)

	// run function based on subcommand
	switch {
	case opt.Config:
		opt.RunConfig()
	case opt.Gen:
		opt.RunGen()
	case opt.Open:
		opt.RunOpen()
	case opt.Fetch:
		opt.RunFetch()
	case opt.Test:
		opt.RunTest()
	case opt.Submit:
		opt.RunSubmit()
	}
	return
}
