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
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// ServerHealth is returned by splunk Server server/health endpoint
type ServerHealth struct {
	Links struct {
	} `json:"links"`
	Origin    string    `json:"origin"`
	Updated   time.Time `json:"updated"`
	Generator struct {
		Build   string `json:"build"`
		Version string `json:"version"`
	} `json:"generator"`
	Entry []struct {
		Name    string    `json:"name"`
		ID      string    `json:"id"`
		Updated time.Time `json:"updated"`
		Links   struct {
			Alternate string `json:"alternate"`
			List      string `json:"list"`
			Details   string `json:"details"`
		} `json:"links"`
		Author string `json:"author"`
		ACL    struct {
			App        string `json:"app"`
			CanList    bool   `json:"can_list"`
			CanWrite   bool   `json:"can_write"`
			Modifiable bool   `json:"modifiable"`
			Owner      string `json:"owner"`
			Perms      struct {
				Read  []string      `json:"read"`
				Write []interface{} `json:"write"`
			} `json:"perms"`
			Removable bool   `json:"removable"`
			Sharing   string `json:"sharing"`
		} `json:"acl"`
		Fields struct {
			Required []interface{} `json:"required"`
			Optional []interface{} `json:"optional"`
			Wildcard []interface{} `json:"wildcard"`
		} `json:"fields"`
		Content struct {
			EaiACL interface{} `json:"eai:acl"`
			Health string      `json:"health"`
		} `json:"content"`
	} `json:"entry"`
	Paging struct {
		Total   int `json:"total"`
		PerPage int `json:"perPage"`
		Offset  int `json:"offset"`
	} `json:"paging"`
	Messages []interface{} `json:"messages"`
}

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

// AppDownload struct for response from splunk server after it downloads the file.
type AppDownload struct {
	Links struct {
		Create string `json:"create"`
		Reload string `json:"_reload"`
	} `json:"links"`
	Origin    string    `json:"origin"`
	Updated   time.Time `json:"updated"`
	Generator struct {
		Build   string `json:"build"`
		Version string `json:"version"`
	} `json:"generator"`
	Entry []struct {
		Name    string    `json:"name"`
		ID      string    `json:"id"`
		Updated time.Time `json:"updated"`
		Links   struct {
			Alternate string `json:"alternate"`
			List      string `json:"list"`
			Reload    string `json:"_reload"`
			Edit      string `json:"edit"`
			Remove    string `json:"remove"`
			Package   string `json:"package"`
		} `json:"links"`
		Author string `json:"author"`
		ACL    struct {
			App            string `json:"app"`
			CanChangePerms bool   `json:"can_change_perms"`
			CanList        bool   `json:"can_list"`
			CanShareApp    bool   `json:"can_share_app"`
			CanShareGlobal bool   `json:"can_share_global"`
			CanShareUser   bool   `json:"can_share_user"`
			CanWrite       bool   `json:"can_write"`
			Modifiable     bool   `json:"modifiable"`
			Owner          string `json:"owner"`
			Perms          struct {
				Read  []string `json:"read"`
				Write []string `json:"write"`
			} `json:"perms"`
			Removable bool   `json:"removable"`
			Sharing   string `json:"sharing"`
		} `json:"acl"`
		Content struct {
			AttributionLink            string      `json:"attribution_link"`
			Author                     string      `json:"author"`
			Build                      int64       `json:"build"`
			CheckForUpdates            bool        `json:"check_for_updates"`
			Configured                 bool        `json:"configured"`
			Core                       bool        `json:"core"`
			Description                string      `json:"description"`
			Details                    string      `json:"details"`
			Disabled                   bool        `json:"disabled"`
			EaiACL                     interface{} `json:"eai:acl"`
			InstallSourceChecksum      string      `json:"install_source_checksum"`
			Label                      string      `json:"label"`
			Location                   string      `json:"location"`
			ManagedByDeploymentClient  bool        `json:"managed_by_deployment_client"`
			Name                       string      `json:"name"`
			ShowInNav                  bool        `json:"show_in_nav"`
			SourceLocation             string      `json:"source_location"`
			StateChangeRequiresRestart bool        `json:"state_change_requires_restart"`
			Status                     string      `json:"status"`
			Version                    string      `json:"version"`
			Visible                    bool        `json:"visible"`
		} `json:"content"`
	} `json:"entry"`
	Paging struct {
		Total   int `json:"total"`
		PerPage int `json:"perPage"`
		Offset  int `json:"offset"`
	} `json:"paging"`
	Messages []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"messages"`
}

// Creds struct contains username and password to auth to Splunkbase
type Creds struct {
	username string
	password string
	splunkp  string
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

// Server struct
type Server struct {
	name string
	port string
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
	if os.Getenv("SPLUNK_PASSWORD") == "" {
		logger().Fatalf("SPLUNK_PASSWORD Not Set: ", os.Getenv("SPLUNK_PASSWORD"))
	}
	cr.username = os.Getenv("SBASE_U")
	cr.password = os.Getenv("SBASE_P")
	cr.splunkp = os.Getenv("SPLUNK_PASSWORD")
	return cr
}

// ServerInfo returns Server struct with info from .env
func ServerInfo() *Server {
	logger().Debug("Start")
	s := new(Server)
	if os.Getenv("SPLUNK_SERVER") == "" {
		s.name = "localhost"
	} else {
		s.name = os.Getenv("SPLUNK_SERVER")
	}
	if os.Getenv("SPLUNK_SERVER_PORT") == "" {
		s.port = "8089"
	} else {
		s.port = os.Getenv("SPLUNK_SERVER_PORT")
	}
	return s
}

// CheckSplunk will see if the splunk server is up and we can make REST API requests to it, if not, it will err.
func CheckSplunk(cr *Creds) (b *ServerHealth, err error) {
	logger().Debug("Start")
	s := ServerInfo()
	b = new(ServerHealth)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := &http.Client{Transport: tr}
	url := "https://" + s.name + ":" + s.port + "/services/server/health/splunkd?output_mode=json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger().Errorln(err)
		return b, err
	}
	req.SetBasicAuth("admin", cr.splunkp)
	res, err := cl.Do(req)
	if err != nil {
		logger().Errorln(err)
		return b, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Errorln(err)
		return b, err
	}
	err = json.Unmarshal(body, &b)
	if err != nil {
		logger().Errorln(err)
		return b, err
	}
	return b, nil

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
func DownloadApp(a *App, c *Creds) (ad *AppDownload, err error) {
	logger().Debug("Start")
	s := ServerInfo()
	ad = new(AppDownload)
	name := "https://splunkbase.splunk.com/app/" + a.UID.String() + "/release/" + a.Release[0].Name + "/download/"
	data := url.Values{}
	data.Set("name", name)
	data.Set("update", "true")
	data.Set("filename", "true")
	data.Set("auth", c.auth)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := &http.Client{Transport: tr}
	url := "https://" + s.name + ":" + s.port + "/services/apps/local/?output_mode=json"
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		logger().Errorln(err)
		return ad, err
	}
	req.SetBasicAuth("admin", c.splunkp)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := cl.Do(req)
	if err != nil {
		logger().Errorln(err)
		return ad, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Errorln(err)
		return ad, err
	}
	err = json.Unmarshal(body, &ad)
	return ad, nil
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

// ZipFile grabs the file from docker and zips it up in the local directory
func ZipFile(ad *AppDownload, c *Creds, path string) (fn string, err error) {
	logger().Debug("Start")
	dirName := ad.Entry[0].Content.Name
	fn = ad.Entry[0].Content.Name + "_" + ad.Entry[0].Content.Version + ".tar.gz"
	logger().Debugln(fn)
	_, err = CheckSplunk(c)
	if err != nil {
		logger().Errorln(err)
		return fn, err
	}
	// Need to hook this into the container creation so we get the container name direct.
	sl := "testing_so1_1:" + ad.Entry[0].Content.SourceLocation
	com := "docker cp " + sl + " " + path + " && " + "tar -zcf " + path + fn + " " + path + dirName + " && " + "rm -rf " + path + dirName
	err = RunCMD("bash", []string{"-c", com})
	if err != nil {
		logger().Errorln(err)
		return fn, err
	}
	return fn, nil
}

// RunCMD runs a command in the shell
func RunCMD(c string, a []string) error {
	logger().Debug("Start")
	cmd := exec.Command(c, a...)
	cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	if err != nil {
		logger().Errorln(err)
		return err
	}
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

//var checks bool
//var quiet bool
var debug bool

func init() {
	flag.StringVar(&appid, "a", "", "AppID to Download")
	//flag.BoolVar(&checks, "c", false, "Perform Checks on Splunk Server")
	//flag.BoolVar(&quiet, "q", false, "Disable Logging Entries and just output json")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if appid == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// if quiet != false {
	// 	log.SetOutput(ioutil.Discard)

	// }
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
	ad, err := DownloadApp(z, Cr)
	if err != nil {
		logger().Fatalln(err)
	}
	// Check if there is already a dir structure for this app, if not, create it.
	filePath := "apps/" + z.UID.String() + "/" + ad.Entry[0].Name + "/"
	path := filePath + ad.Entry[0].Content.Version + "/"
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
	s.Suffix = "  > Getting files from Splunk"
	_, err = ZipFile(ad, Cr, path)
	if err != nil {
		logger().Fatalln(err)
	}
	s.FinalMSG = "Download for " + z.UID.String() + " - " + z.Title + " Version: " + ad.Entry[0].Content.Version + " is complete!\nFiles located at " + path + " "
	s.Stop()
	//fmt.Printf("Done! %v : %v Downloaded: Detailed Information in %vappinfo.json", z.UID.String(), z.Title, filePath)

}
