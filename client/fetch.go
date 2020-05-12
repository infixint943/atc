package cln

import (
	cfg "atc/config"
	pkg "atc/packages"

	"bytes"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FindCountdown parses countdown (if exists) from timer board
func FindCountdown(contest string, link url.URL) (int64, error) {
	// Instantiate http client to make requests
	// link already points to desired url link
	c := cfg.Session.Client
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return 0, err
	} else if len(body) == 0 {
		// contest page doesn't exist
		err := fmt.Errorf("contest %v doesn't exist", contest)
		return 0, err
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find(".fixtime").First().Text()
	// extract time remaining from start time string.
	// may go negative (if contest has already started)
	dur, err := time.Parse("2006-01-02 15:04:05-0700", val)
	secs := int64(dur.Sub(time.Now()).Seconds())
	return secs, nil
}

// StartCountdown starts countdown of dur seconds
func StartCountdown(dur int64) {
	// run timer till it runs out
	pkg.LiveUI.Start()
	for ; dur > 0; dur-- {
		h := fmt.Sprintf("%d:", dur/(60*60))
		m := fmt.Sprintf("0%d:", (dur/60)%60)
		s := fmt.Sprintf("0%d", dur%60)
		pkg.LiveUI.Print(h + m[len(m)-3:] + s[len(s)-2:])
		time.Sleep(time.Second)
	}
	// remove timer data from screen
	pkg.LiveUI.Print()
	return
}

// FetchProbs finds all problems present in the contest
func FetchProbs(contest string, link url.URL) ([]string, error) {
	// extract details from tasks dashboard (instead of main page, idk why tho)
	c := cfg.Session.Client
	link.Path = path.Join(link.Path, "tasks")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	}

	var probs []string
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find("table tbody tr").Each(func(_ int, s *goquery.Selection) {
		prob := strings.TrimSpace(s.Find("td").First().Text())
		probs = append(probs, strings.ToLower(prob))
	})
	return probs, nil
}

// FetchTests extracts test cases of the problem(s) in contest
// Returns 2d slice mapping to input and output
func FetchTests(contest string, link url.URL) ([][]string, [][]string, error) {

	c := cfg.Session.Client
	link.Path = path.Join(link.Path, "tasks_print")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, nil, err
	}

	// splInp will hold input of each problem
	// splOut maps to splInp with the output data
	var splInp, splOut [][]string
	// Iterate over every problem
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".col-sm-12").Not(".next-page").Find(".lang-en").Each(func(_ int, prob *goquery.Selection) {
		var inp, out []string

		prob.Find(".part > section > pre").Each(func(i int, spl *goquery.Selection) {
			if i == 0 {
				// skip over input format
				return
			}
			if i%2 == 1 {
				inp = append(inp, spl.Text())
			} else {
				out = append(out, spl.Text())
			}
		})

		splInp = append(splInp, inp)
		splOut = append(splOut, out)
	})
	return splInp, splOut, nil
}
