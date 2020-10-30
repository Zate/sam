package pkgr

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

const (
	filePath = "../sbase/apps/"
)

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
	flag.StringVar(&appid, "a", "", "AppID to Repackage")
	flag.BoolVar(&debug, "d", false, "Turn on Debug")
	flag.Parse()
	if appid == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.SetFormatter(&log.JSONFormatter{})
	logger().Debug("Init")
	if debug != false {
		log.SetLevel(log.DebugLevel)
	}
	logger().Debug("Init")
}

func main() {
	logger().Debug("Start")
	svrTypes := []string{"idx", "fwd", "shc"}
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Color("red", "bold")
	s.Start()
	s.Suffix = "  > Parsing appinfo.json for AppID " + appid
	app := appinfo(appid)
	s.Suffix = "  > Unpacking " + app.Appid + "Version: " + app.LatestVersion
	UnpackApp(&app, svrTypes)
	s.Suffix = "  > Parsing files in apps/" + fmt.Sprint(app.UID) + "/{idx | shc | fws}"
	err := parsePackages(svrTypes, &app)
	if err != nil {
		logger().Fatalln(err)
	}
	logger().Debug(len(app.Packages[0].Objects))
	s.FinalMSG = "  > Completed Unpacking and Parsing Objects in {idx | shc | fws}"
	s.Stop()

}
