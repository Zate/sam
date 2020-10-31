package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Creds struct contains username and password to auth to Splunkbase
type Creds struct {
	username string
	password string
	auth     string
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
	logger().Debug("Start")
	cr := new(Creds)
	if os.Getenv("SBASE_U") == "" {
		logger().Fatalf("SBASE_U Not Set: %v", os.Getenv("SBASE_U"))
		os.Exit(1)
	}
	if os.Getenv("SBASE_P") == "" {
		logger().Fatalf("SBASE_P Not Set: ", os.Getenv("SBASE_P"))
		os.Exit(1)
	}
	cr.username = os.Getenv("SBASE_U")
	cr.password = os.Getenv("SBASE_P")
	return cr
}

// AuthToken function
func AuthToken(cr *Creds) string {
	logger().Debug("Start")
	formData := url.Values{
		"username": {cr.username},
		"password": {cr.password},
	}
	res, err := http.PostForm("https://splunkbase.splunk.com/api/account:login/", formData)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	var auth Feed
	err = xml.Unmarshal(body, &auth)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	return auth.ID
}
