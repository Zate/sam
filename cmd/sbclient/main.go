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
	"github.com/zate/sam/pkg/sbase"
)

var (
	appid       string
	debug       bool
	manifest    string
	catalogpath string
	catalog     bool
	dlpath      string
	dlpathInfo  os.FileInfo
)

func init() {
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.StringVar(&manifest, "m", sbase.FilePath+"manifest.json", "Path to manifest.json")
	flag.StringVar(&catalogpath, "cp", sbase.FilePath+"catalog.json", "Path to catalog.json")
	flag.StringVar(&dlpath, "p", "", "Path to extract the addon")
	flag.BoolVar(&catalog, "c", false, "Update Catalog (Not Functioning)")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if (appid == "" && catalog == false) || catalog == true {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if dlpath != "" {
		// lets validate this path exists
		dlpathInfo, err := sbase.CheckPath(dlpath)
		if err != nil {
			fmt.Printf("%v \n\nUsage:\n", err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		sbase.Logger().Debugln(dlpathInfo.Name())
	}
	if err := godotenv.Load(); err != nil {
		sbase.Logger().Debug("No .env file found")
	}
	log.SetFormatter(&log.JSONFormatter{})
	sbase.Logger().Debug("Init")
	if debug != false {
		log.SetLevel(log.DebugLevel)
	}
	m := sbase.FilePath + "manifest.json"
	if manifest != "" {
		m = manifest
	}
	sbase.LoadManifest(m)
}

func main() {
	sbase.Logger().Debug("Start")
	// if catalog != false {

	// 	getAllApps(catalogpath)

	// 	os.Exit(1)

	// }
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("red", "bold")
	s.Start()
	s.Suffix = "  > Loading Environment"
	Cr := sbase.LoadCreds()
	s.Suffix = "  > Getting Splunkbase Auth Token"
	Cr.Auth = sbase.AuthToken(Cr)
	s.Suffix = "  > Parsing AppInfo"
	z := sbase.GetApp(appid)
	z.LatestVersion = z.Release[0].Name
	z.LatestRelease = string(z.Release[0].ID)
	s.Suffix = "  > Downloading " + appid
	tgz := sbase.DownloadApp(&z, Cr)
	fPath := sbase.FilePath + fmt.Sprint(z.UID) + "/" + z.Appid + "/"
	path := fPath + z.LatestVersion + "/"
	s.Suffix = "  > Creating " + path
	err := sbase.CheckDir(path)
	if err != nil {
		sbase.Logger().Fatalln(err)
	}
	s.Suffix = "  > Writing appinfo.json"
	file, err := json.MarshalIndent(z, "", " ")
	if err != nil {
		sbase.Logger().Fatalln(err)
	}
	err = ioutil.WriteFile(sbase.FilePath+fmt.Sprint(z.UID)+"/appinfo.json", file, 0644)
	sbase.Logger().Debug("Writing " + sbase.FilePath + fmt.Sprint(z.UID) + "/appinfo.json")
	if err != nil {
		sbase.Logger().Fatalln(err)

	}
	err = ioutil.WriteFile(path+z.Appid+"_"+z.LatestVersion+".tar.gz", tgz, 0644)
	sbase.Logger().Debug("Writing " + path + z.Appid + "_" + z.LatestVersion + ".tar.gz")
	if err != nil {
		sbase.Logger().Fatalln(err)
		os.Exit(1)
	}
	svrTypes := []string{"idx", "fwd", "shc"}
	s.Suffix = "  > Unpacking " + z.Appid + "Version: " + z.LatestVersion
	sbase.UnpackApp(&z, svrTypes, dlpath)
	sbase.Logger().Debug(len(z.Packages[0].Objects))
	s.FinalMSG = "Download for " + fmt.Sprint(z.UID) + " - " + z.Title + " Version: " + z.LatestVersion + " is complete!\nFiles located at " + path + " "
	s.Stop()
}
