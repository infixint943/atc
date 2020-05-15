package cmd

import (
	cln "atc/client"
	cfg "atc/config"
	pkg "atc/packages"
)

// RunSubmit is called on running cf submit
func (opt Opts) RunSubmit() {
	// check if problem id is present
	if opt.problem == "" {
		pkg.Log.Error("No problem id found")
		return
	}
	// find code file to submit
	file, err := cln.FindSourceFiles(opt.File)
	pkg.PrintError(err, "Failed to select source file")
	// find template config to use
	t, err := cln.FindTmpltsConfig(file)
	pkg.PrintError(err, "Failed to select template configuration")

	// check login status
	usr, err := cln.LoggedInUsr()
	pkg.PrintError(err, "Failed to check login status")
	if usr == "" {
		// exit if no saved login configurations found
		if cfg.Session.Handle == "" || cfg.Session.Passwd == "" {
			pkg.Log.Error("No login details configured")
			pkg.Log.Notice("Configure login details through cf config")
			return
		}
		// attempt relogin
		pkg.Log.Warning("No logged in user session found")
		pkg.Log.Info("Attempting relogin: " + cfg.Session.Handle)
		status, err := cln.Relogin()
		pkg.PrintError(err, "Failed to login")
		if status == true {
			// logged in successfully
			pkg.Log.Success("Login successful")
		} else {
			pkg.Log.Error("Login failed")
			pkg.Log.Notice("Configure login details through 'cf config'")
			return
		}
	} else {
		// output handle details of current user
		// this is in else loop, since current user is already
		// being displayed during relogin above
		pkg.Log.Notice("Current user: " + usr)
	}

	// main submit code runs here
	err = cln.Submit(opt.contest, opt.problem, t.LangID, file, opt.link)
	pkg.PrintError(err, "Failed to submit source code")
	pkg.Log.Success("Submitted")
	// watch submission verdict
	// opt.watch()
}

// func (opt Opts) watch() {
// 	// infinite loop till verdicts declared
// 	pkg.LiveUI.Start()
// 	for {
// 		// fetch submission status from contest every second
// 		start := time.Now()
//
// 		data, err :=
// 	}
// }
