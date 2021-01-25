package sbase

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	// FilePath refers to where we store the apps.
	FilePath = "apps/"
)

// Appinfo takes an id (string) and reads in the appinfo.json for that app, returning an App struct.
// func Appinfo(id string) (a App) {
// 	defer TimeTrack(time.Now())
// 	Logger().Debug("Start")
// 	info := FilePath + id + "/appinfo.json"
// 	Logger().Debug(info)

// 	f, err := ReadFile(info)
// 	ChkErr(err)

// 	err = json.Unmarshal([]byte(f), &a)
// 	ChkErr(err)
// 	return a
// }

// GetURL returns Body from Splunkbase API request.
func GetURL(u string) (body []byte) {
	defer TimeTrack(time.Now())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u, nil)
	ChkErr(err)
	res, err := cl.Do(req)

	ChkErr(err)
	body, err = ioutil.ReadAll(res.Body)

	ChkErr(err)
	return body
}

// GetApp takes an appid and calls Splunkbase API to get info about the app.
func (sb *SBase) GetApp(a string) (app App) {
	defer TimeTrack(time.Now())
	// var app App
	u := "https://splunkbase.splunk.com/api/v1/app/" + a + "/"
	b := GetURL(u)
	err := json.Unmarshal(b, &app)
	ChkErr(err)
	u = u + "release/"
	b = GetURL(u)
	err = json.Unmarshal(b, &app.Release)
	ChkErr(err)
	return app
}

// // CheckVer verifies that the version you are requesting is the latest one for download on the website.
// func (sb *SBase) CheckVer(a string) string {
// 	return ""
// }

// DownloadApp func downloads the app through Splunk server
func (sb *SBase) DownloadApp(a *App) (body []byte) {
	defer TimeTrack(time.Now())
	url := "https://splunkbase.splunk.com/app/" + fmt.Sprint(a.UID) + "/release/" + a.LatestVersion + "/download/"
	// fmt.Println(url)
	cl := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	ChkErr(err)
	req.Header.Add("X-Auth-Token", sb.Creds.Auth)
	res, err := cl.Do(req)
	// fmt.Printf("%v\n", fmt.Sprint(res.StatusCode))
	ChkErr(err)
	body, err = ioutil.ReadAll(res.Body)
	// fmt.Printf("%v\n", fmt.Sprint(string(body)))
	ChkErr(err)
	return body
}

// CheckDir looks to see if the directory structure for this app exists, if not, creates it.
func CheckDir(path string) error {
	defer TimeTrack(time.Now())
	// Apps / AppID / Name / Ver
	if err := os.MkdirAll(path, 0755); os.IsExist(err) {
		Logger().Errorf("%v already exists", path)
		return nil
	}
	return nil
}
