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
	appspath    string
	catalogpath string
	catalog     bool
	dlpath      string
	dlpathInfo  os.FileInfo
	s           = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sb          = sbase.New()
)

func init() {
	defer sbase.TimeTrack(time.Now())
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.StringVar(&appspath, "m", sbase.FilePath, "Path to download apps")
	flag.StringVar(&catalogpath, "cp", sbase.FilePath+"catalog.json", "Path to catalog.json")
	flag.StringVar(&dlpath, "p", "", "Path to extract the addon")
	flag.BoolVar(&catalog, "c", false, "Update Catalog (Not Functioning)")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if (appid == "" && catalog == false) || catalog == true {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if debug != false {
		log.SetLevel(log.DebugLevel)
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
	sbase.FilePath = appspath
	sb.LoadManifest(appspath)
}

func main() {
	defer sbase.TimeTrack(time.Now())
	// if catalog != false {

	// 	getAllApps(catalogpath)

	// 	os.Exit(1)

	// }
	if debug == false {
		s.Color("red", "bold")
		s.Start()
		s.Suffix = "  > Loading Environment"
	}
	sb.LoadCreds()
	if debug == false {
		s.Suffix = "  > Getting Splunkbase Auth Token"
	}
	sb.AuthToken()
	if debug == false {
		s.Suffix = "  > Parsing AppInfo"
	}
	z := sb.GetApp(appid)
	z.LatestVersion = z.Release[0].Name
	z.LatestRelease = fmt.Sprint(z.Release[0].ID)
	if debug == false {
		s.Suffix = "  > Downloading " + appid
	}
	tgz := sb.DownloadApp(&z)
	fPath := sbase.FilePath + fmt.Sprint(z.UID) + "/" + z.Appid + "/"
	path := fPath + z.LatestVersion + "/"
	if debug == false {
		s.Suffix = "  > Creating " + path
	}
	err := sbase.CheckDir(path)
	sbase.ChkErr(err)
	if debug == false {
		s.Suffix = "  > Writing appinfo.json"
	}
	file, err := json.MarshalIndent(z, "", " ")
	sbase.ChkErr(err)
	err = ioutil.WriteFile(sbase.FilePath+fmt.Sprint(z.UID)+"/appinfo.json", file, 0644)
	sbase.ChkErr(err)
	err = ioutil.WriteFile(path+z.Appid+"_"+z.LatestVersion+".tar.gz", tgz, 0644)
	sbase.ChkErr(err)
	svrTypes := []string{"idx", "fwd", "shc"}
	if debug == false {
		s.Suffix = "  > Unpacking " + z.Appid + "Version: " + z.LatestVersion
	}
	sbase.UnpackApp(&z, svrTypes, dlpath)
	if debug == false {
		s.FinalMSG = "Download for " + fmt.Sprint(z.UID) + " - " + z.Title + " Version: " + z.LatestVersion + " is complete!\nFiles located at " + path + " \n\n"
		s.Stop()
	}
}
