package cmd

import (
	cln "atc/client"
	cfg "atc/config"
	pkg "atc/packages"

	"github.com/AlecAivazis/survey/v2"
)

// RunConfig is called on running atc config
func (opt Opts) RunConfig() {
	var choice int

	err := survey.AskOne(&survey.Select{
		Message: "Select configuration:",
		Options: []string{
			"Login to atcoder",
			"Add new code template",
			"Remove code template",
			"Other misc preferences",
		},
	}, &choice, survey.WithValidator(survey.Required))
	pkg.PrintError(err, "")

	switch choice {
	case 0:
		login()
		// case 1:
		// 	addTmplt()
		// case 2:
		// 	remTmplt()
		// case 3:
		// 	miscPrefs()
	}
	return
}

func login() {
	// check if logged in user exists
	if cfg.Session.Handle != "" {
		pkg.Log.Success("Current user: " + cfg.Session.Handle)
		pkg.Log.Warning("Current session will be overwritten")
	}
	// take input of username / password
	creds := struct{ Usr, Passwd string }{}
	err := survey.Ask([]*survey.Question{
		{
			Name:     "usr",
			Prompt:   &survey.Input{Message: "Username:"},
			Validate: survey.Required,
		}, {
			Name:     "passwd",
			Prompt:   &survey.Password{Message: "Password:"},
			Validate: survey.Required,
		},
	}, &creds)
	pkg.PrintError(err, "")
	// login and check login status
	pkg.Log.Info("Logging in")
	flag, err := cln.Login(creds.Usr, creds.Passwd)
	pkg.PrintError(err, "Login failed")
	// login was successful
	if flag == true {
		pkg.Log.Success("Login successful")
		pkg.Log.Notice("Welcome " + cfg.Session.Handle)
	} else {
		// login failed
		pkg.Log.Error("Login failed.")
		pkg.Log.Notice("Check credentials and retry")
	}
	return
}
