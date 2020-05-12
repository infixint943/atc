package pkg

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func parseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetReqBody executes a GET request to url and returns the request body
func GetReqBody(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode == 404 {
		// return blank page, to signify contest doesn't exist
		return nil, nil
	}
	return parseBody(resp)
}

// PostReqBody executes a POST request (with values: data) to url and returns the request body
func PostReqBody(client *http.Client, url string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	return parseBody(resp)
}

// FindHandle scrapes handle from REQUEST body
func FindHandle(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	sel := doc.Find(".navbar-right .dropdown").First().Next()
	val := GetText(sel, ".dropdown-toggle")
	return val
}

// FindCsrf extracts Csrf from REQUEST body
func FindCsrf(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := GetAttr(doc.Selection, "[name=\"csrf_token\"]", "value")
	return val
}

// RedirectCheck prevents redirection and returns requested page info
func RedirectCheck(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
