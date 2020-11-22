package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// DList is a struct that has the list of urls to download
type DList struct {
	URLs []string
}

// D var to make a url list
var D DList

const allURL = "https://splunkbase.splunk.com/api/v1/app/?include=all"

//var Catalog *os.File

func init() {

}

func main() {
	urls := make(chan string, 4)
	offset := 0
	limit := 10
	total := getTotal(allURL + "&offset=0&limit=1")
	logger().Debugf("Items is: %v %v %v", total, offset, limit)
	// s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// s.Color("red", "bold")
	// s.Start()
	var wg sync.WaitGroup
	go func() {
		for {
			u := allURL + "&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
			fmt.Println("Adding " + u)
			urls <- u
			offset = offset + limit
			if (offset) > total {
				logger().Debug("Exit URL Generator")
				return
			}
		}
	}()

	go func(wg *sync.WaitGroup) {
		for {

			defer wg.Done()
			url := <-urls
			fmt.Println(url)
			// s.Suffix = "  > DownLoading " + url
			// res := GetURL(url)

			// fH, err := os.OpenFile("catalog.json", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
			// if err != nil {
			// 	logger().Fatalln(err)
			// 	os.Exit(1)
			// }
			// writer := bufio.NewWriter(fH)
			// defer fH.Close()

			// fmt.Fprintln(writer, string(res))
			// writer.Flush()

		}
	}(&wg)
	wg.Wait()
	// s.FinalMSG = "  > Download Done"
	// s.Stop()
}

func getTotal(s string) (t int) {

	var allapps AllApps
	res := GetURL(s)
	err := json.Unmarshal(res, &allapps)
	if err != nil {
		logger().Debug(s)
		logger().Fatalln(err)
		os.Exit(1)
	}
	return allapps.Total
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
		os.Exit(1)
	}
	res, err := cl.Do(req)

	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	return body
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

// doExist checks if the object s exists
func doExist(s string) bool {
	_, err := os.Stat(s)
	if err != nil {
		return false
	}
	return true
}

// Catalog is a struct containing all the Results to search
type Catalog struct {
	Results []Result  `json:"results"`
	Updated time.Time `json:"updated"`
}

// AllApps struct contains all the Apps
type AllApps struct {
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Total   int      `json:"total"`
	Results []Result `json:"results"`
}

// Result struct with everything about the result
type Result struct {
	UID                      int       `json:"uid"`
	Appid                    string    `json:"appid"`
	Title                    string    `json:"title"`
	CreatedTime              time.Time `json:"created_time"`
	PublishedTime            time.Time `json:"published_time"`
	UpdatedTime              time.Time `json:"updated_time"`
	LicenseName              string    `json:"license_name"`
	Type                     string    `json:"type"`
	LicenseURL               string    `json:"license_url"`
	Description              string    `json:"description"`
	Access                   string    `json:"access"`
	AppinspectPassed         bool      `json:"appinspect_passed"`
	Path                     string    `json:"path"`
	InstallMethodDistributed string    `json:"install_method_distributed"`
	InstallMethodSingle      string    `json:"install_method_single"`
	Support                  string    `json:"support"`
	CreatedBy                struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
	} `json:"created_by"`
	DownloadCount DownloadCountWrapper `json:"download_count"`
	InstallCount  InstallCountWrapper  `json:"install_count"`
	Categories    []string             `json:"categories"`
	Icon          string               `json:"icon"`
	Screenshots   []string             `json:"screenshots"`
	Rating        struct {
		Count   int `json:"count"`
		Average int `json:"average"`
	} `json:"rating"`
	Release           interface{}   `json:"release"`
	Releases          []interface{} `json:"releases"`
	Documentation     string        `json:"documentation"`
	Manifest          interface{}   `json:"manifest"`
	ArchiveStatus     string        `json:"archive_status"`
	IsArchived        bool          `json:"is_archived"`
	FedrampValidation string        `json:"fedramp_validation"`
}

// DownloadCountWrapper to fix bad API's
type DownloadCountWrapper struct {
	DownloadCount int
}

// InstallCountWrapper to fix bad API's
type InstallCountWrapper struct {
	InstallCount int
}

func notmain() {
	// page through all of the available info using
	// offset=<number>
	// limit=<number>
	// var goRoutine sync.WaitGroup

	// workers := 8
	// offsetChan := make(chan int, workers-1)
	// goRoutine.Add(1)
	// logger().Debugf("Total is %v", total)
	// go offsets(offsetChan, total)
	// var Catalog Catalog

	// go func(offsetChan chan int) {
	// 	for {
	// 		goRoutine.Add(1)
	// 		url := allURL + "&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	// 		var allapps AllApps
	// 		res := GetURL(url)
	// 		err := json.Unmarshal(res, &allapps)

	// 		if err != nil {
	// 			logger().Debug(url)
	// 			logger().Fatalln(err)
	// 			os.Exit(1)
	// 		}
	// 		total = allapps.Total
	// 		limit = allapps.Limit
	// 		offset = <-offsetChan
	// 		Catalog.Results = append(Catalog.Results, allapps.Results...)
	// 		logger().Debugf("Offset: %d	Limit: %d	Total: %d", offset, limit, total)
	// 		goRoutine.Done()
	// 		if (offset) > total {
	// 			break
	// 		}
	// 	}
	// }(offsetChan)
	// goRoutine.Wait()

	// Catalog.updateNow("catalog.json")
	return

}

func offsets(o chan int, t int) {

	offset := 0
	limit := 10
	total := t

	for {
		o <- offset
		offset = offset + limit
		if (offset) > total {
			break
		}
	}
}

func (m *Catalog) updateNow(c string) {
	m.Updated = time.Now()
	mout, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		logger().Errorln(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(c, mout, 0644)
	if err != nil {
		logger().Errorln(err)
		os.Exit(1)
	}
	return

}

// UnmarshalJSON is to unf*ck someones bad API.
func (w *DownloadCountWrapper) UnmarshalJSON(data []byte) (err error) {

	if dc, err := strconv.Atoi(string(data)); err == nil {
		w.DownloadCount = dc
		return nil
	}
	w.DownloadCount = 0
	return nil
}

// UnmarshalJSON is to unf*ck someones bad API.
func (w *InstallCountWrapper) UnmarshalJSON(data []byte) (err error) {

	if dc, err := strconv.Atoi(string(data)); err == nil {
		w.InstallCount = dc
		return nil
	}
	w.InstallCount = 0
	return nil
}
