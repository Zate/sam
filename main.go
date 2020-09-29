package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

// Creds struct contains username and password to auth to Splunkbase
type Creds struct {
	username string
	password string
}

// Auth struct contains artifacts gotten from splunk Okta auth when successful
type Auth struct {
	statusCode string
	status     string
	message    string
	sid        string
	ssoid      string
}

// SplunkAuthURL is the Splunk Okta Auth endpoint.
var SplunkAuthURL = "https://account.splunk.com/api/v1/okta/auth"

// GetApp will download this appid from splunkbase
// https://splunkbase.splunk.com/app/2919/release/5.4.0/download/
// https://splunkbase.splunk.com/api/v1/app/2919
// might need to crawl https://splunkbase.splunk.com/app/2919 to grab version info
//https://docs.splunk.com/Documentation/SplunkbaseAPI/current/SBAPI/About

// GetApp will download this appid from splunkbase
func GetApp(aid string) string {
	cr := LoadCreds()
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	a := &http.Client{Jar: jar}
	res, err := a.PostForm(SplunkAuthURL, url.Values{
		"password": {cr.password},
		"username": {cr.username},
	})
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(SplunkAuthURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	log.WithFields(log.Fields{
		"func":    "GetApp",
		"app_id":  aid,
		"cookies": res.Cookies(),
	}).Print("Cookies After Auth")

	surl := "https://splunkbase.splunk.com/app/" + aid
	// Request the HTML page.

	res, err = a.Get(surl)
	if err != nil {
		log.Fatal(err)
	}
	// bdy, err := ioutil.ReadAll(res.Body)
	// fmt.Print(string(bdy))

	defer res.Body.Close()

	u, err = url.Parse(surl)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	log.WithFields(log.Fields{
		"func":         "GetApp",
		"app_id":       aid,
		"cookies":      res.Cookies(),
		"sent_cookies": res.Request.Cookies(),
	}).Print("Cookies After Get Splunkbase URL")

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// <sb-release-select-click u-for="download-modal" sb-selector="release-version" sb-target="5.4.0" class="u.hide@sm u.item:6/12@xl u.item:1/1@lg u.item:1/1@md u.btn:green" sb-href="/app/2919/release/5.4.0/download/" data-ol-has-click-handler="">
	//                                                 Download
	//                                             </sb-release-select-click>

	// Find the review items
	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		// log.Print(i)
		id, ex := s.Attr("id")
		if ex && id == "release-version" {
			ver, _ := s.Find("option").First().Attr("value")
			// log.Print(ver)
			dlurl := surl + "/release/" + ver + "/download/"
			// log.Print(dlurl)
			res, err = a.Get(dlurl)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			u, err = url.Parse(dlurl)
			if err != nil {
				log.Fatal(err)
			}
			for _, cookie := range jar.Cookies(u) {
				fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
			}
			log.WithFields(log.Fields{
				"func":             "GetApp",
				"app_id":           aid,
				"cookies":          res.Cookies(),
				"sent_cookies":     res.Request.Cookies(),
				"Response_Headers": res.Header,
			}).Print("Cookies After Download")

			out, err := os.Create(path.Base(res.Request.URL.String()))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// Write the body to file
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}
			_, err = io.Copy(out, res.Body)

		}
		// fmt.Println(i, title)

	})

	return aid

}

// SplunkAuth takes a set of creds, and auths to Splunk Okta endpoint, returns Auth struct
func SplunkAuth() (client http.Client) {

	cr := LoadCreds()
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client = http.Client{Jar: jar}
	resp, err := client.PostForm(SplunkAuthURL, url.Values{
		"password": {cr.password},
		"username": {cr.username},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print(resp.Cookies())

	return client
}

// LoadCreds grabs username and password from SB_USER and SB_PASSWD env vars
func LoadCreds() (c *Creds) {
	// Load from env.
	cr := new(Creds)

	if os.Getenv("SB_USER") == "" {
		log.Print(os.Getenv("SB_USER"))
		log.Fatal("SB_USER Not Set")
	}

	if os.Getenv("SB_PASSWD") == "" {
		log.Print(os.Getenv("SB_PASSWD"))
		log.Fatal("SB_PASSWD Not Set")

	}

	cr.username = os.Getenv("SB_USER")
	log.Print("SB_USER read")
	cr.password = os.Getenv("SB_PASSWD")
	log.Print("SB_PASSWD read")

	return cr
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Print("Init done")
	// If debug is set - SAM_DEBUG=1 log.SetLevel(log.DebugLevel)
	if os.Getenv("SAM_DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	}

}

func main() {
	log.Print("Main Started")
	log.Debug("Debug is on")
	GetApp("2919")
}
