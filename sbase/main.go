package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	appid    string
	debug    bool
	manifest string
)

const (
	filePath = "apps/"
)

func init() {
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.StringVar(&manifest, "m", filePath+"manifest.json", "Path to manifest.json")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if appid == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err := godotenv.Load(); err != nil {
		logger().Debug("No .env file found")
	}
	log.SetFormatter(&log.JSONFormatter{})
	logger().Debug("Init")
	if debug != false {
		log.SetLevel(log.DebugLevel)
	}
	m := filePath + "manifest.json"
	if manifest != "" {
		m = manifest
	}
	LoadManifest(m)
}

func main() {
	logger().Debug("Start")
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("red", "bold")
	s.Start()
	s.Suffix = "  > Loading Environment"
	Cr := LoadCreds()
	s.Suffix = "  > Getting Splunkbase Auth Token"
	Cr.auth = AuthToken(Cr)
	s.Suffix = "  > Parsing AppInfo"
	z := GetApp(appid)
	z.LatestVersion = z.Release[0].Name
	z.LatestRelease = string(z.Release[0].ID)
	s.Suffix = "  > Downloading " + appid
	tgz := DownloadApp(&z, Cr)
	fPath := filePath + fmt.Sprint(z.UID) + "/" + z.Appid + "/"
	path := fPath + z.LatestVersion + "/"
	s.Suffix = "  > Creating " + path
	err := CheckDir(path)
	if err != nil {
		logger().Fatalln(err)
	}
	s.Suffix = "  > Writing appinfo.json"
	file, err := json.MarshalIndent(z, "", " ")
	if err != nil {
		logger().Fatalln(err)
	}
	err = ioutil.WriteFile(filePath+fmt.Sprint(z.UID)+"/appinfo.json", file, 0644)
	logger().Debug("Writing " + filePath + fmt.Sprint(z.UID) + "/appinfo.json")
	if err != nil {
		logger().Fatalln(err)

	}
	err = ioutil.WriteFile(path+z.Appid+"_"+z.LatestVersion+".tar.gz", tgz, 0644)
	logger().Debug("Writing " + path + z.Appid + "_" + z.LatestVersion + ".tar.gz")
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	svrTypes := []string{"idx", "fwd", "shc"}
	s.Suffix = "  > Unpacking " + z.Appid + "Version: " + z.LatestVersion
	UnpackApp(&z, svrTypes)
	logger().Debug(len(z.Packages[0].Objects))
	s.FinalMSG = "Download for " + fmt.Sprint(z.UID) + " - " + z.Title + " Version: " + z.LatestVersion + " is complete!\nFiles located at " + path + " "
	s.Stop()
}
