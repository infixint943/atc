package cln

import (
	cfg "atc/config"
	pkg "atc/packages"

	"encoding/hex"
	"net/url"
	"path"

	"github.com/infixint943/cookiejar"
	"github.com/oleiade/serrure/aes"
)

// Login tries logginging in with user creds
func Login(usr, passwd string) (bool, error) {
	// instantiate http client, but remove
	// past user sessions to prevent redirection
	jar, _ := cookiejar.New(nil)
	c := cfg.Session.Client
	c.Jar = jar

	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, "login")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return false, err
	}

	// Hidden form data
	csrf := pkg.FindCsrf(body)

	// Post form (aka login using creds)
	body, err = pkg.PostReqBody(&c, link.String(), url.Values{
		"csrf_token": {csrf},
		"username":   {usr},
		"password":   {passwd},
	})
	if err != nil {
		return false, err
	}

	usr = pkg.FindHandle(body)
	if usr != "" {
		// create aes 256 encryption and encode as
		// hex string and save to sessions.json
		enc, _ := aes.NewAES256Encrypter(usr, nil)
		ed, _ := enc.Encrypt([]byte(passwd))
		ciphertext := hex.EncodeToString(ed)
		// update sessions data
		cfg.Session.Cookies = jar
		cfg.Session.Handle = usr
		cfg.Session.Passwd = ciphertext
		cfg.SaveSession()
	}
	return (usr != ""), nil
}

// LoggedInUsr checks and returns whether
// current session is logged in
func LoggedInUsr() (string, error) {
	// fetch home page and check if logged in
	c := cfg.Session.Client
	link, _ := url.Parse(cfg.Settings.Host)
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return "", err
	}

	return pkg.FindHandle(body), nil
}
