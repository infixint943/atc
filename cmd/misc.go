package cmd

import (
	cfg "atc/config"

	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Version of current binary tool
const Version = "0.9.0"

type (
	// Opts is struct docopt binds flag data to
	Opts struct {
		Config bool `docopt:"config"`
		Gen    bool `docopt:"gen"`
		Open   bool `docopt:"open"`
		Fetch  bool `docopt:"fetch"`
		Test   bool `docopt:"test"`
		Submit bool `docopt:"submit"`

		Info []string `docopt:"<info>"`

		All    bool   `docopt:"--all"`
		File   string `docopt:"--file"`
		IgCase bool   `docopt:"--ignore-case"`
		Exp    int    `docopt:"--ignore-exp"`
		Tl     int    `docopt:"--time-limit"`
		Custom bool   `docopt:"--custom"`

		contest string
		problem string
		dirPath string
		link    url.URL
	}

	// Env are global (generic and non-generic) variables
	Env struct {
		// generic variables
		handle string `env:"${handle}"`
		date   string `env:"${date}"`
		time   string `env:"${time}"`

		// non-generic variables
		Contest string `env:"${contest}"`
		Problem string `env:"${problem}"`
		Idx     string `env:"${idx}"`
		File    string `env:"${file}"`
	}
)

// FindContestData extracts contest / problem id from path
func (opt *Opts) FindContestData() {
	// path to current directory
	currPath, _ := os.Getwd()

	if len(opt.Info) == 0 {
		// no contest id given in flags. Fetch from folder path
		data := strings.Split(currPath, string(os.PathSeparator))
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10

		// cleans path to return dir path to root folder
		clean := func(i int) string {
			str := filepath.Join(data[i:]...)
			return strings.TrimSuffix(currPath, str)
		}
		// find last directory matching 'Settings.WSName'
		for i := sz - 1; i >= 0; i-- {
			// current folder name matches configured WSName
			if data[i] == cfg.Settings.WSName {
				// set contestClass, contest and problem
				opt.contest = data[i+1]
				opt.problem = data[i+2]
				currPath = clean(i)
				break
			}
		}
	} else if _, err := url.ParseRequestURI(opt.Info[0]); err == nil {
		// url given in the flags. parse data from url
		data := strings.Split(opt.Info[0], "/")
		// prevent out-of-bounds accessing
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10
		// iterate over each part of url and
		// find first part matching criteria
		for i := 0; i < sz; i++ {
			if data[i] == "contests" {
				opt.contest = data[i+1]
				// convert '_' to '-' in problem
				data[i+3] = strings.ReplaceAll(data[i+3], "_", "-")
				// remove contest prefix from problem string
				opt.problem = strings.TrimPrefix(data[i+3], data[i+1]+"-")
				break
			}
		}
	} else {
		// parse from command line args (for example, abc123 e)
		data := append(opt.Info, make([]string, 10)...)
		// set contest and problem id from flags
		opt.contest = data[0]
		opt.problem = data[1]
	}
	// convert problem id to lowercase
	opt.problem = strings.ToLower(opt.problem)
	// set path to folder containing contClass
	opt.dirPath = filepath.Join(currPath, cfg.Settings.WSName)
	// set common link to contest
	// dereference the url variable
	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, "contests", opt.contest)
	opt.link = *link

	return
}

// ReplPlaceholder replaces all global variables in text
// with their respective values. Non-generic are passed as map
func (e Env) ReplPlaceholder(text string) string {
	// set handle/date/time
	e.handle = cfg.Session.Handle
	e.date = time.Now().Format("02-01-06")
	e.time = time.Now().Format("15:04:05")

	// replace string data
	repl := func(old, new string) string {
		return strings.ReplaceAll(text, old, new)
	}
	// omit ${idx} = 0
	if e.Idx == "0" {
		e.Idx = ""
	}

	// iterate over struct and replace variables
	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("env")
		val := v.Field(i).String()
		text = repl(tag, val)
	}

	return text
}

/*
Parsing structure of problems
-----------------------------
- WSName
  - ${contest}
    - ${problem}
*/
