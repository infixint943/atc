package cln

import (
	cfg "atc/config"
	pkg "atc/packages"

	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
)

// Submit uploads form data and submits user code
func Submit(contest, problem, langID, file string, link url.URL) error {

	c := cfg.Session.Client
	link.Path = path.Join(link.Path, "submit")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("contest %v doesn't exist", contest)
		return err
	}

	// read source file
	data, _ := ioutil.ReadFile(file)
	// hidden form data
	csrf := pkg.FindCsrf(body)
	// merge problem string with contest string
	problem = strings.ReplaceAll(contest+"_"+problem, "-", "_")
	// post form data
	body, err = pkg.PostReqBody(&c, link.String(), url.Values{
		"data.TaskScreenName": {problem},
		"data.LanguageId":     {langID},
		"sourceCode":          {string(data)},
		"csrf_token":          {csrf},
	})
	if err != nil {
		return err
	}
	return nil
}
