package cmd

import (
	pkg "atc/packages"

	"os/exec"
	"path"
	"runtime"
	"strings"
)

// RunOpen is called on running `cf open`
func (opt Opts) RunOpen() {
	// check if contest id is present
	if opt.contest == "" {
		pkg.Log.Error("No contest id found")
		return
	}
	link := opt.link
	// open problems page (all problems)
	if opt.problem == "" {
		link.Path = path.Join(link.Path, "tasks_print")
	} else {
		problem := strings.ReplaceAll(opt.contest+"_"+opt.problem, "-", "_")
		link.Path = path.Join(link.Path, "tasks", problem)
	}
	// open page in default browser
	browserOpen(link.String())
	return
}

func browserOpen(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
	return
}
