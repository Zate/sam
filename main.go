package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

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
	// Get Addon
	// Package App
	// Package Addon
	// Combine
	// Deploy

}
