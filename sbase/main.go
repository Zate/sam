package main

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

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
func CheckSplunk(cr *Creds) error {
	s := ServerInfo()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := "https://" + s.name + ":" + s.port + "/services/server/info?output_mode=json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	req.SetBasicAuth("admin", cr.splunkp)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(body))

	return nil

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

func init() {
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
	err := CheckSplunk(Cr)
	if err != nil {
		log.Fatalln(err)
	}

}
