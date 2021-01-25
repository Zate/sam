package sbase

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// LoadCreds grabs username and password from SB_USER and SB_PASSWD env vars
func (sb *SBase) LoadCreds() {
	defer TimeTrack(time.Now())
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
	sb.Creds = cr
	return
}

// AuthToken function
func (sb *SBase) AuthToken() {
	defer TimeTrack(time.Now())
	formData := url.Values{
		"username": {sb.Creds.Username},
		"password": {sb.Creds.Password},
	}
	res, err := http.PostForm("https://splunkbase.splunk.com/api/account:login/", formData)
	ChkErr(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	ChkErr(err)
	var auth Feed
	//fmt.Printf("%v", string(body))
	err = xml.Unmarshal(body, &auth)
	ChkErr(err)
	sb.Creds.Auth = auth.ID
	return
}
