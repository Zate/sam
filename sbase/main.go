package main

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// App stuff
type App struct {
	UID                      json.Number `json:"uid"`
	Appid                    string      `json:"appid"`
	Title                    string      `json:"title"`
	CreatedTime              time.Time   `json:"created_time"`
	PublishedTime            time.Time   `json:"published_time"`
	UpdatedTime              time.Time   `json:"updated_time"`
	LicenseName              string      `json:"license_name"`
	Type                     string      `json:"type"`
	LicenseURL               string      `json:"license_url"`
	Description              string      `json:"description"`
	Access                   string      `json:"access"`
	AppinspectPassed         bool        `json:"appinspect_passed"`
	Path                     string      `json:"path"`
	InstallMethodDistributed string      `json:"install_method_distributed"`
	InstallMethodSingle      string      `json:"install_method_single"`
	DownloadCount            int         `json:"download_count"`
	InstallCount             int         `json:"install_count"`
	ArchiveStatus            string      `json:"archive_status"`
	IsArchived               bool        `json:"is_archived"`
	FedrampValidation        string      `json:"fedramp_validation"`
	Release                  Release     `json:"release"`
}

// Release struct
type Release []struct {
	ID                        json.Number   `json:"id"`
	App                       json.Number   `json:"app"`
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
	}
	if os.Getenv("SBASE_P") == "" {
		logger().Fatalf("SBASE_P Not Set: ", os.Getenv("SBASE_P"))
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
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var auth Feed
	xml.Unmarshal(body, &auth)
	return auth.ID
}

// GetURL returns Body from Splunkbase API request.
func GetURL(u string) (body []byte) {
	logger().Debug("Start")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		logger().Fatalln(err)
	}
	res, err := cl.Do(req)
	if err != nil {
		logger().Fatalln(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Fatalln(err)
	}
	return body
}

// GetApp takes an appid and calls Splunkbase API to get info about the app.
func GetApp(a string) (app *App) {
	logger().Debug("Start")
	app = new(App)
	u := "https://splunkbase.splunk.com/api/v1/app/" + a + "/"
	b := GetURL(u)
	err := json.Unmarshal(b, &app)
	if err != nil {
		logger().Fatalln(err)
	}
	u = u + "release/"
	b = GetURL(u)
	err = json.Unmarshal(b, &app.Release)
	if err != nil {
		logger().Fatalln(err)
	}
	return app
}

// DownloadApp func downloads the app through Splunk server
func DownloadApp(a *App, c *Creds) (body []byte, err error) {
	logger().Debug("Start")
	url := "https://splunkbase.splunk.com/app/" + a.UID.String() + "/release/" + a.Release[0].Name + "/download/"
	cl := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger().Errorln(err)
		return body, err
	}
	req.Header.Add("X-Auth-Token", c.auth)
	res, err := cl.Do(req)
	if err != nil {
		logger().Errorln(err)
		return body, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Errorln(err)
		return body, err
	}
	return body, nil
}

// CheckDir looks to see if the directory structure for this app exists, if not, creates it.
func CheckDir(path string) error {
	logger().Debug("Start")
	// Apps / AppID / Name / Ver
	logger().Debugln(path)
	if err := os.MkdirAll(path, 0755); os.IsExist(err) {
		logger().Errorf("%v already exists", path)
		return nil
	}
	logger().Debugf("%v created", path)
	return nil
}

func logger() *log.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return log.WithField("file", filename).WithField("function", fn)
}

var appid string
var debug bool

func init() {
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if appid == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	log.SetFormatter(&log.JSONFormatter{})
	logger().Debug("Init")
	if debug != false {
		log.SetLevel(log.DebugLevel)
	}
	// Probably want to put something here to read in the list of appid's from a yml file or something.
	logger().Debug("Init")
}

func main() {
	logger().Debug("Start")
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Color("red", "bold")
	s.FinalMSG = "Complete!\nNew line!\nAnother one!\n"
	s.Start()
	s.Suffix = "  > Loading Environment" // Start the spinner
	Cr := LoadCreds()
	s.Suffix = "  > Getting Splunkbase Auth Token"
	Cr.auth = AuthToken(Cr)
	s.Suffix = "  > Parsing AppInfo"
	z := GetApp(appid)
	s.Suffix = "  > Downloading " + appid
	tgz, err := DownloadApp(z, Cr)
	if err != nil {
		logger().Fatalln(err)
	}
	// Check if there is already a dir structure for this app, if not, create it.
	filePath := "apps/" + z.UID.String() + "/" + z.Appid + "/"
	path := filePath + z.Release[0].Name + "/"
	s.Suffix = "  > Creating " + path
	err = CheckDir(path)
	if err != nil {
		logger().Fatalln(err)
	}
	s.Suffix = "  > Writing appinfo.json"
	file, err := json.MarshalIndent(z, "", " ")
	if err != nil {
		logger().Fatalln(err)
	}
	err = ioutil.WriteFile(filePath+"appinfo.json", file, 0644)
	if err != nil {
		logger().Fatalln(err)
	}
	err = ioutil.WriteFile(path+z.Appid+"_"+z.Release[0].Name+".tar.gz", tgz, 0644)
	if err != nil {
		logger().Fatalln(err)
	}
	s.FinalMSG = "Download for " + z.UID.String() + " - " + z.Title + " Version: " + z.Release[0].Name + " is complete!\nFiles located at " + path + " "
	s.Stop()
}
