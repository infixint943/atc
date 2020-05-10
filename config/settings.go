package cfg

import (
	pkg "atc/packages"

	"encoding/json"
	"io/ioutil"
	"os"
)

// Settings holds configured settings data of the tool
var Settings struct {
	DfltTmplt  int    `json:"default_template"`
	GenOnFetch bool   `json:"gen_on_fetch"`
	Host       string `json:"host"`
	Proxy      string `json:"proxy"`
	WSName     string `json:"workspace_name"`
}

func init() {
	Settings.DfltTmplt = -1
	Settings.GenOnFetch = false
	Settings.Host = "https://atcoder.jp"
	Settings.Proxy = ""
	Settings.WSName = "atcoder"
}

var settPath string

// InitSettings reads settings.json file
func InitSettings(path string) {
	// set settings.json file path
	settPath = path

	file, err := ioutil.ReadFile(settPath)
	if err != nil {
		pkg.Log.Warning("File settings.json doesn't exist")
		pkg.Log.Info("Creating settings.json file")
		SaveSettings()
	}
	json.Unmarshal(file, &Settings)
}

// SaveSettings to settings.json file
func SaveSettings() {
	file, err := os.Create(settPath)
	pkg.PrintError(err, "Failed to create settings.json file")

	body, _ := json.MarshalIndent(Settings, "", "\t")
	file.Write(body)
}
