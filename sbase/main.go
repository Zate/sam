package main

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

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
	// {
	// 	"links": {},
	// 	"origin": "https://192.168.0.16:8089/services/server/health",
	// 	"updated": "2020-10-23T19:44:37+00:00",
	// 	"generator": {
	// 	  "build": "f57c09e87251",
	// 	  "version": "8.1.0"
	// 	},
	// 	"entry": [
	// 	  {
	// 		"name": "splunkd",
	// 		"id": "https://192.168.0.16:8089/services/server/health/splunkd",
	// 		"updated": "1970-01-01T00:00:00+00:00",
	// 		"links": {
	// 		  "alternate": "/services/server/health/splunkd",
	// 		  "list": "/services/server/health/splunkd",
	// 		  "details": "/services/server/health/splunkd/details"
	// 		},
	// 		"author": "system",
	// 		"acl": {
	// 		  "app": "",
	// 		  "can_list": true,
	// 		  "can_write": true,
	// 		  "modifiable": false,
	// 		  "owner": "system",
	// 		  "perms": {
	// 			"read": [
	// 			  "admin",
	// 			  "splunk-system-role"
	// 			],
	// 			"write": []
	// 		  },
	// 		  "removable": false,
	// 		  "sharing": "system"
	// 		},
	// 		"fields": {
	// 		  "required": [],
	// 		  "optional": [],
	// 		  "wildcard": []
	// 		},
	// 		"content": {
	// 		  "eai:acl": null,
	// 		  "health": "yellow"
	// 		}
	// 	  }
	// 	],
	// 	"paging": {
	// 	  "total": 1,
	// 	  "perPage": 30,
	// 	  "offset": 0
	// 	},
	// 	"messages": []
	//   }

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

// VWrap wraps Splunk Version to get around mixed string / int types.
// type VWrap json.Number

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

// LatestRelease struct for sotring info from latest release
// type LatestRelease struct {
// 	App            json.Number   `json:"app"`
// 	Name           string        `json:"name"`
// 	ReleaseNotes   string        `json:"release_notes"`
// 	SplunkVersions []string      `json:"splunk_versions"`
// 	CIMVersions    []interface{} `json:"CIM_versions"`
// 	Public         bool          `json:"public"`
// }

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

	// <?xml version="1.0" encoding="utf-8"?>
	// 	<feed xmlns="http://www.w3.org/2005/Atom">
	//     <title>Authentication Token</title>
	//     <updated>2020-10-21T21:23:40.332535+00:00</updated>
	//     <id>ogkzudn4pyxcbsplphroo50zdogigjzn</id>
	// </feed>
}

// Server struct
type Server struct {
	name string
	port string
	// tls bool
}

// // App is a json struct containing all the info about an app.
// type App struct {
// 	UID                      int    `json:"uid"`
// 	Appid                    string `json:"appid"`
// 	Title                    string `json:"title"`
// 	CreatedTime              string `json:"created_time"`
// 	PublishedTime            string `json:"published_time"`
// 	UpdatedTime              string `json:"updated_time"`
// 	LicenseName              string `json:"license_name"`
// 	AppType                  string `json:"type"`
// 	LicenseURL               string `json:"license_url"`
// 	Description              string `json:"description"`
// 	Access                   string `json:"access"`
// 	AppInspectPassed         bool   `json:"appinspect_passed"`
// 	Path                     string `json:"path"`
// 	InstallMethodDistributed string `json:"install_method_distributed"`
// 	InstallMethodSingle      string `json:"install_method_single"`
// 	DownloadCount            int    `json:"download_count"`
// 	InstallCount             int    `json:"install_count"`
// 	ArchiveStatus            string `json:"archive_status"`
// 	IsArchived               bool   `json:"is_archived"`
// 	FedrampValidation        string `json:"fedramp_validation"`
// }

// UnmarshalJSON overirde for Splunk Version stuff.
// func (w *VWrap) UnmarshalJSON(data []byte) (err error) {
// 	if ver, err := strconv.Atoi(string(data)); err == nil {
// 		str := strconv.Itoa(ver)
// 		*w = VWrap(str)
// 		return nil
// 	}
// 	var str string
// 	err = json.Unmarshal(data, &str)
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal([]byte(str), w)
// }

// LoadCreds grabs username and password from SB_USER and SB_PASSWD env vars
func LoadCreds() (c *Creds) {
	// Load from env.
	cr := new(Creds)

	if os.Getenv("SBASE_U") == "" {
		log.Print(os.Getenv("SBASE_U"))
		log.Fatal("SBASE_U Not Set")
	}

	if os.Getenv("SBASE_P") == "" {
		log.Print(os.Getenv("SBASE_P"))
		log.Fatal("SBASE_P Not Set")

	}

	if os.Getenv("SPLUNK_PASSWORD") == "" {
		log.Print(os.Getenv("SPLUNK_PASSWORD"))
		log.Fatal("SPLUNK_PASSWORD Not Set")
	}

	cr.username = os.Getenv("SBASE_U")
	cr.password = os.Getenv("SBASE_P")
	cr.splunkp = os.Getenv("SPLUNK_PASSWORD")

	return cr
}

// ServerInfo returns Server struct with info from .env
func ServerInfo() *Server {
	// Check for env stuff loaded.
	// SPLUNK_SERVER
	// SPLUNK_SERVER_PORT
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
	s := ServerInfo()
	b = new(ServerHealth)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := "https://" + s.name + ":" + s.port + "/services/server/health/splunkd?output_mode=json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return b, err
	}
	req.SetBasicAuth("admin", cr.splunkp)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return b, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return b, err
	}
	//fmt.Println(string(body))

	err = json.Unmarshal(body, &b)

	if err != nil {
		log.Println(err)
		return b, err
	}

	return b, nil

}

// AuthToken function
func AuthToken(cr *Creds) string {
	//AUTH=`curl -sS -d "username=${SBASE_U}&password=${SBASE_P}" -X POST https://splunkbase.splunk.com/api/account:login/ | grep -o -P '(?<=<id>).*(?=</id>)'`

	// cr := LoadCreds()
	formData := url.Values{
		"username": {cr.username},
		"password": {cr.password},
	}

	res, err := http.PostForm("https://splunkbase.splunk.com/api/account:login/", formData)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var auth Feed

	xml.Unmarshal(body, &auth)
	// fmt.Println(auth.Title)

	return auth.ID

}

// GetURL returns Body from Splunkbase API request.
func GetURL(u string) (body []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u, nil)

	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(body))

	return body
}

// GetApp takes an appid and calls Splunkbase API to get info about the app.
func GetApp(a string) (app *App) {
	app = new(App)
	u := "https://splunkbase.splunk.com/api/v1/app/" + a + "/"
	b := GetURL(u)
	err := json.Unmarshal(b, &app)
	if err != nil {
		log.Fatalln(err)
	}

	// rel := new(Release)
	u = u + "release/"
	b = GetURL(u)

	err = json.Unmarshal(b, &app.Release)
	if err != nil {
		log.Fatalln(err)
	}

	// app.Release = *rel

	// lrel := new(LatestRelease)
	// u = u + string(app.Release[0].ID) + "/"
	// fmt.Println(u)
	// b = GetURL(u)
	// fmt.Println(string(b))
	// err = json.Unmarshal(b, &app.Release[0])
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// app.Release[0].App = lrel.App
	// app.Release[0].Name = lrel.Name
	// app.Release[0].ReleaseNotes = lrel.ReleaseNotes
	// app.Release[0].SplunkVersions = lrel.SplunkVersions
	// app.Release[0].CIMVersions = lrel.CIMVersions
	// app.Release[0].Public = lrel.Public

	return app
}

var appid string
var checks bool
var quiet bool

func init() {
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.BoolVar(&checks, "c", false, "Perform Checks on Splunk Server")
	flag.BoolVar(&quiet, "q", false, "Disable Logging Entries and just output json")
	flag.Parse()
	if quiet != false {
		log.SetOutput(ioutil.Discard)

	}
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	log.SetFormatter(&log.JSONFormatter{})

	log.Print("Init done")
	// If debug is set - SAM_DEBUG=1 log.SetLevel(log.DebugLevel)
	if os.Getenv("SAM_DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	}

}

func main() {

	Cr := LoadCreds()
	Cr.auth = AuthToken(Cr)
	// log.Println(Cr.auth)
	if checks == true {
		b, err := CheckSplunk(Cr)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Server Health is: %v", b.Entry[0].Content.Health)
		if b.Entry[0].Content.Health == "" || b.Entry[0].Content.Health == "red" {
			log.Fatalf("Server doesnt exist or is red: %v", b.Entry[0].Content.Health)
		}
	}
	// Get a list of all Splunkbase apps.  Maybe we cache this locally to browse through at some stage?
	// Get AppInfo based on Appid
	if appid != "" {
		z := GetApp(appid)
		log.Println(z.Release[0].Name)
		if quiet != false {
			tmp, err := json.Marshal(z)

			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(tmp))
			// v := string(z.Release[0].Name)
			// fmt.Printf("Version: %v", v)
		}

	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

}
