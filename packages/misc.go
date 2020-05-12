package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

// Global Variables for different UI formatting
var (
	writer = os.Stderr
	Green  = color.New(color.FgGreen)
	Blue   = color.New(color.FgBlue)
	Red    = color.New(color.FgRed)
	Yellow = color.New(color.FgYellow)

	Log struct {
		Success, Notice, Info, Error,
		Warning func(text ...interface{})
	}
	// LiveUI to print live data to terminal
	LiveUI struct {
		count int
		isAPI bool
		Start func()
		Print func(text ...string)
	}
)

func init() {
	// Initialise colored text output
	Log.Success = func(text ...interface{}) { Green.Fprintln(writer, text...) }
	Log.Notice = func(text ...interface{}) { fmt.Fprintln(writer, text...) }
	Log.Info = func(text ...interface{}) { Blue.Fprintln(writer, text...) }
	Log.Error = func(text ...interface{}) { Red.Fprintln(writer, text...) }
	Log.Warning = func(text ...interface{}) { Yellow.Fprintln(writer, text...) }

	// Initialise Live rendering output
	LiveUI.isAPI = false
	LiveUI.Start = func() { LiveUI.count = 0 }
	LiveUI.Print = func(text ...string) {
		// clear last count lines from terminal
		for i := 0; !LiveUI.isAPI && i < LiveUI.count; i++ {
			ansi.CursorPreviousLine(1)
			ansi.EraseInLine(2)
		}
		// count number of lines in text
		LiveUI.count = 1
		for _, str := range text {
			LiveUI.count += strings.Count(str, "\n")
			fmt.Println(str)
		}
	}
}

// // IsAPI configures output settings based on flag
// func IsAPI(flag bool) {
// 	// disable color
// 	color.NoColor = flag
// 	// disable line overwriting
// 	LiveUI.isAPI = flag
// }

// PrintError outputs error (with custom message)
// and exits the program execution (if err != nil)
func PrintError(err error, desc string) {
	if err != nil {
		if desc != "" {
			Log.Error(desc)
		}
		Log.Error(err.Error())
		os.Exit(0)
	}
}

// CreateFile copies data to dst (create if not exists)
// Returns absolute path to destination file
func CreateFile(data, dst string) string {
	out, err := os.Create(dst)
	PrintError(err, "File "+dst+" couldn't be created!")
	defer out.Close()

	out.WriteString(data)
	return dst
}

// GetText extracts text from particular html data
func GetText(sel *goquery.Selection, query string) string {
	str := sel.Find(query).Text()
	return strings.TrimSpace(str)
}

// GetAttr extracts attribute valur of particular html data
func GetAttr(sel *goquery.Selection, query, attr string) string {
	str := sel.Find(query).AttrOr(attr, "")
	return strings.TrimSpace(str)
}
