package sbase

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Creds struct contains username and password to auth to Splunkbase
type Creds struct {
	Username string
	Password string
	Auth     string
}

// Auth Strict
type Auth struct {
	XMLName xml.Name `xml:"xml"`
	Feed    Feed     `xml:"feed"`
}

// Feed Struct
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	ID      string   `xml:"id"`
}

// LoadCreds grabs username and password from SB_USER and SB_PASSWD env vars
func LoadCreds() (c *Creds) {
	Logger().Debug("Start")
	cr := new(Creds)
	if os.Getenv("SBASE_U") == "" {
		Logger().Fatalf("SBASE_U Not Set: %v", os.Getenv("SBASE_U"))
		os.Exit(1)
	}
	if os.Getenv("SBASE_P") == "" {
		Logger().Fatalf("SBASE_P Not Set: ", os.Getenv("SBASE_P"))
		os.Exit(1)
	}
	cr.Username = os.Getenv("SBASE_U")
	cr.Password = os.Getenv("SBASE_P")
	return cr
}

// AuthToken function
func AuthToken(cr *Creds) string {
	Logger().Debug("Start")
	formData := url.Values{
		"username": {cr.Username},
		"password": {cr.Password},
	}
	res, err := http.PostForm("https://splunkbase.splunk.com/api/account:login/", formData)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	var auth Feed
	err = xml.Unmarshal(body, &auth)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	return auth.ID
}
