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
	appid    string
	debug    bool
	quiet    bool
	owrite   bool
	appspath string
	// appver      string
	// catalogpath string
	catalog    bool
	dlpath     string
	dlpathInfo os.FileInfo
	s          = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sb         = sbase.New()
)

func init() {
	defer sbase.TimeTrack(time.Now())
	flag.StringVar(&appid, "a", "", "AppID to Download")
	flag.StringVar(&appspath, "m", sbase.FilePath, "Path to download apps")
	// flag.StringVar(&catalogpath, "cp", sbase.FilePath+"catalog.json", "Path to catalog.json")
	flag.StringVar(&dlpath, "p", "", "Path to extract the package")
	// flag.StringVar(&appver, "v", "", "Specific Version string of the app to get")
	flag.BoolVar(&catalog, "c", false, "Update Catalog (Not Functioning)")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.BoolVar(&quiet, "q", false, "Silence All Output")
	flag.BoolVar(&owrite, "o", false, "Overwite output path")
	flag.Parse()
	if appid == "" || catalog == true {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if debug != false {
		log.SetLevel(log.DebugLevel)
	}
	if err := godotenv.Load(); err != nil {
		sbase.Logger().Debug("No .env file found")
	}
	log.SetFormatter(&log.JSONFormatter{})
	if dlpath != "" {
		if sbase.DoExist(dlpath) == true {
			if owrite == false {
				sbase.Logger().Errorf("%v exists and -o (overwrite) not specified", dlpath)
				os.Exit(1)
			}
			err := os.RemoveAll(dlpath)
			sbase.ChkErr(err)
		}
		// lets validate this path exists
		dlpathInfo, err := sbase.CheckPath(dlpath)
		if err != nil {
			fmt.Printf("%v \n\nUsage:\n", err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		if debug == true {
			sbase.Logger().Debugln(dlpathInfo.Name())
		}
	} else {
		sbase.FilePath = appspath
		_ = sbase.CheckDir(appspath)
		//sb.LoadManifest(appspath)
	}
}

func main() {
	defer sbase.TimeTrack(time.Now())
	if debug == false && quiet == false {
		s.Color("red", "bold")
		s.Start()
		s.Suffix = "  > Loading Environment"
	}
	sb.LoadCreds()
	if debug == false && quiet == false {
		s.Suffix = "  > Getting Splunkbase Auth Token"
	}
	sb.AuthToken()
	if debug == false && quiet == false {
		s.Suffix = "  > Parsing AppInfo"
	}
	z := sb.GetApp(appid)
	if debug == true {
		d, err := json.MarshalIndent(z, "", "  ")
		sbase.ChkErr(err)
		fmt.Printf("%s\n", d)
	}
	for a, x := range z.Release {
		if x.Public == true {
			z.LatestRelease = fmt.Sprint(z.Release[a].ID)
			z.LatestVersion = z.Release[a].Name
			if debug == true {
				sbase.Logger().Printf("%v is %v - ver: %v", fmt.Sprint(z.UID), z.LatestRelease, z.LatestVersion)
			}
			break
		}
	}
	//sb.Catalog[z.Appid] = &z
	//var ma *map[string]string = new(map[string]string)
	//sb.Manifest.Apps = ma
	//ma[]
	//ma[z.Appid] = fmt.Sprint(sb.Catalog[z.Appid].UID)
	sb.Manifest.Apps = make(map[int]string)
	sb.LoadManifest(appspath)
	sb.Manifest.Apps[z.UID] = z.Appid
	//fmt.Sprint(sb.Catalog[z.Appid].UID)
	//fmt.Sprint(z.UID)

	if debug == false && quiet == false {
		s.Suffix = "  > Downloading " + appid
	}
	if dlpath != "" {
		sb.DLOnly(dlpath, &z)
		if debug == false && quiet == false {
			s.FinalMSG = "Download for " + fmt.Sprint(z.UID) + " - " + z.Title + " Version: " + z.LatestVersion + " is complete!\nFiles located at " + dlpath + " \n\n"
			s.Stop()
		}
		os.Exit(0)
	}
	tgz := sb.DownloadApp(&z)
	fPath := sbase.FilePath + fmt.Sprint(z.UID) + "/" + z.Appid + "/"
	path := fPath + z.LatestVersion + "/"
	if debug == false && quiet == false {
		s.Suffix = "  > Creating " + path
	}
	err := sbase.CheckDir(path)
	sbase.ChkErr(err)
	if debug == false && quiet == false {
		s.Suffix = "  > Writing appinfo.json"
	}
	file, err := json.MarshalIndent(z, "", " ")
	sbase.ChkErr(err)
	err = ioutil.WriteFile(sbase.FilePath+fmt.Sprint(z.UID)+"/appinfo.json", file, 0644)
	sbase.ChkErr(err)
	err = ioutil.WriteFile(path+z.Appid+"_"+z.LatestVersion+".tar.gz", tgz, 0644)
	sbase.ChkErr(err)
	svrTypes := []string{"cm", "ds", "shd"}
	if debug == false && quiet == false {
		s.Suffix = "  > Unpacking " + z.Appid + " Version: " + z.LatestVersion
	}
	sbase.UnpackApp(&z, svrTypes, dlpath)
	sb.UpdateNow(sbase.FilePath)
	if debug == false && quiet == false {
		s.FinalMSG = "Download for " + fmt.Sprint(z.UID) + " - " + z.Title + " Version: " + z.LatestVersion + " is complete!\nFiles located at " + path + " \n\n"
		s.Stop()
	}
}
