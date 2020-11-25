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

// App stuff
// Todo : get this App stuff shared correctly between the packages, maybe its a package on its own?
type App struct {
	UID                      int       `json:"uid,omitempty"`
	Appid                    string    `json:"appid,omitempty"`
	Title                    string    `json:"title,omitempty"`
	CreatedTime              time.Time `json:"created_time,omitempty"`
	PublishedTime            time.Time `json:"published_time,omitempty"`
	UpdatedTime              time.Time `json:"updated_time,omitempty"`
	LicenseName              string    `json:"license_name,omitempty"`
	Type                     string    `json:"type,omitempty"`
	LicenseURL               string    `json:"license_url,omitempty"`
	Description              string    `json:"description,omitempty"`
	Access                   string    `json:"access,omitempty"`
	AppinspectPassed         bool      `json:"appinspect_passed,omitempty"`
	Path                     string    `json:"path,omitempty"`
	InstallMethodDistributed string    `json:"install_method_distributed,omitempty"`
	InstallMethodSingle      string    `json:"install_method_single,omitempty"`
	DownloadCount            int       `json:"download_count,omitempty"`
	InstallCount             int       `json:"install_count,omitempty"`
	ArchiveStatus            string    `json:"archive_status,omitempty"`
	IsArchived               bool      `json:"is_archived,omitempty"`
	FedrampValidation        string    `json:"fedramp_validation,omitempty"`
	LatestVersion            string    `json:"latest_version,omitempty"`
	LatestRelease            string    `json:"latest_release,omitempty"`
	Release                  Release   `json:"release,omitempty"`
	Packages                 []Package `json:"packages,omitempty"`
}

// Release struct
type Release []struct {
	ID                        int           `json:"id"`
	App                       int           `json:"app"`
	Name                      string        `json:"name"`
	ReleaseNotes              string        `json:"release_notes"`
	CIMVersions               []interface{} `json:"CIM_versions"`
	SplunkVersions            []int         `json:"splunk_versions"`
	Public                    bool          `json:"public"`
	PublicEverTrue            bool          `json:"public_ever_true"`
	CreatedDatetime           time.Time     `json:"created_datetime"`
	PublishedDatetime         time.Time     `json:"published_datetime"`
	Size                      int           `json:"size"`
	Filename                  string        `json:"filename"`
	Platform                  string        `json:"platform"`
	IsBundle                  bool          `json:"is_bundle"`
	HasUI                     bool          `json:"has_ui"`
	Approved                  bool          `json:"approved"`
	AppinspectStatus          bool          `json:"appinspect_status"`
	InstallMethodSingle       string        `json:"install_method_single"`
	InstallMethodDistributed  string        `json:"install_method_distributed"`
	RequiresCloudVetting      bool          `json:"requires_cloud_vetting"`
	AppinspectRequestID       interface{}   `json:"appinspect_request_id"`
	CloudVettingRequestID     string        `json:"cloud_vetting_request_id"`
	Python3Acceptance         bool          `json:"python3_acceptance"`
	Python3AcceptanceDatetime time.Time     `json:"python3_acceptance_datetime"`
	Python3AcceptanceUser     int           `json:"python3_acceptance_user"`
	FedrampValidation         string        `json:"fedramp_validation"`
	CloudCompatible           bool          `json:"cloud_compatible"`
}

// Package struct represents a single collection of FSOBjects that make up a Splunk App
type Package struct {
	DType   string     `json:"d_type,omitempty"`
	Objects []FSObject `json:"objects,omitempty"`
}

// FSObject struct represents a file system object
type FSObject struct {
	ID           int         `json:"id,omitempty"`
	RelativePath string      `json:"relative_path,omitempty"`
	Name         string      `json:"name,omitempty"`
	Type         string      `json:"type,omitempty"`
	FileInfo     os.FileInfo `json:"file_info,omitempty"`
}

func appinfo(id string) (a App) {
	Logger().Debug("Start")
	info := FilePath + id + "/appinfo.json"
	Logger().Debug(info)

	f, err := ioutil.ReadFile(info)
	if err != nil {
		Logger().Errorln(err)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(f), &a)
	if err != nil {
		Logger().Errorln(err)
	}
	return a
}

// GetURL returns Body from Splunkbase API request.
func GetURL(u string) (body []byte) {
	Logger().Debug("Start")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	res, err := cl.Do(req)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	return body
}

// GetApp takes an appid and calls Splunkbase API to get info about the app.
func GetApp(a string) (app App) {
	Logger().Debug("Start")
	// var app App
	u := "https://splunkbase.splunk.com/api/v1/app/" + a + "/"
	b := GetURL(u)
	err := json.Unmarshal(b, &app)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	u = u + "release/"
	b = GetURL(u)
	err = json.Unmarshal(b, &app.Release)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	return app
}

// DownloadApp func downloads the app through Splunk server
func DownloadApp(a *App, c *Creds) (body []byte) {
	Logger().Debug("Start")
	url := "https://splunkbase.splunk.com/app/" + fmt.Sprint(a.UID) + "/release/" + a.Release[0].Name + "/download/"
	cl := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	req.Header.Add("X-Auth-Token", c.Auth)
	res, err := cl.Do(req)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	return body
}

// CheckDir looks to see if the directory structure for this app exists, if not, creates it.
func CheckDir(path string) error {
	Logger().Debug("Start")
	// Apps / AppID / Name / Ver
	Logger().Debugln(path)
	if err := os.MkdirAll(path, 0755); os.IsExist(err) {
		Logger().Errorf("%v already exists", path)
		return nil
	}
	Logger().Debugf("%v created", path)
	return nil
}
