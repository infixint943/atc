package cmd

import (
	cfg "atc/config"

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

		Info []string `docopt:"<info>"`

		All bool `docopt:"--all"`

		contest string
		problem string
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
