package main

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"
// )

// // Catalog is a struct containing all the Results to search
// type Catalog struct {
// 	Results []Result  `json:"results"`
// 	Updated time.Time `json:"updated"`
// }

// // AllApps struct contains all the Apps
// type AllApps struct {
// 	Offset  int      `json:"offset"`
// 	Limit   int      `json:"limit"`
// 	Total   int      `json:"total"`
// 	Results []Result `json:"results"`
// }

// // Result struct with everything about the result
// type Result struct {
// 	UID                      int       `json:"uid"`
// 	Appid                    string    `json:"appid"`
// 	Title                    string    `json:"title"`
// 	CreatedTime              time.Time `json:"created_time"`
// 	PublishedTime            time.Time `json:"published_time"`
// 	UpdatedTime              time.Time `json:"updated_time"`
// 	LicenseName              string    `json:"license_name"`
// 	Type                     string    `json:"type"`
// 	LicenseURL               string    `json:"license_url"`
// 	Description              string    `json:"description"`
// 	Access                   string    `json:"access"`
// 	AppinspectPassed         bool      `json:"appinspect_passed"`
// 	Path                     string    `json:"path"`
// 	InstallMethodDistributed string    `json:"install_method_distributed"`
// 	InstallMethodSingle      string    `json:"install_method_single"`
// 	Support                  string    `json:"support"`
// 	CreatedBy                struct {
// 		Username    string `json:"username"`
// 		DisplayName string `json:"display_name"`
// 	} `json:"created_by"`
// 	DownloadCount DownloadCountWrapper `json:"download_count"`
// 	InstallCount  InstallCountWrapper  `json:"install_count"`
// 	Categories    []string             `json:"categories"`
// 	Icon          string               `json:"icon"`
// 	Screenshots   []string             `json:"screenshots"`
// 	Rating        struct {
// 		Count   int `json:"count"`
// 		Average int `json:"average"`
// 	} `json:"rating"`
// 	Release           interface{}   `json:"release"`
// 	Releases          []interface{} `json:"releases"`
// 	Documentation     string        `json:"documentation"`
// 	Manifest          interface{}   `json:"manifest"`
// 	ArchiveStatus     string        `json:"archive_status"`
// 	IsArchived        bool          `json:"is_archived"`
// 	FedrampValidation string        `json:"fedramp_validation"`
// }

// // DownloadCountWrapper to fix bad API's
// type DownloadCountWrapper struct {
// 	DownloadCount int
// }

// // InstallCountWrapper to fix bad API's
// type InstallCountWrapper struct {
// 	InstallCount int
// }

// const allURL = "https://splunkbase.splunk.com/api/v1/app/?include=all"

// func main() {
// 	// page through all of the available info using
// 	// offset=<number>
// 	// limit=<number>
// 	// var goRoutine sync.WaitGroup
// 	offset := 0
// 	limit := 10
// 	total := getTotal(allURL + "&offset=0&limit=1")
// 	fmt.Printf("Items is: %v %v %v", total, offset, limit)
// 	// workers := 8
// 	// offsetChan := make(chan int, workers-1)
// 	// goRoutine.Add(1)
// 	// logger().Debugf("Total is %v", total)
// 	// go offsets(offsetChan, total)
// 	// var Catalog Catalog

// 	// go func(offsetChan chan int) {
// 	// 	for {
// 	// 		goRoutine.Add(1)
// 	// 		url := allURL + "&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
// 	// 		var allapps AllApps
// 	// 		res := GetURL(url)
// 	// 		err := json.Unmarshal(res, &allapps)

// 	// 		if err != nil {
// 	// 			logger().Debug(url)
// 	// 			logger().Fatalln(err)
// 	// 			os.Exit(1)
// 	// 		}
// 	// 		total = allapps.Total
// 	// 		limit = allapps.Limit
// 	// 		offset = <-offsetChan
// 	// 		Catalog.Results = append(Catalog.Results, allapps.Results...)
// 	// 		logger().Debugf("Offset: %d	Limit: %d	Total: %d", offset, limit, total)
// 	// 		goRoutine.Done()
// 	// 		if (offset) > total {
// 	// 			break
// 	// 		}
// 	// 	}
// 	// }(offsetChan)
// 	// goRoutine.Wait()

// 	// Catalog.updateNow("catalog.json")
// 	return

// }

// func offsets(o chan int, t int) {

// 	offset := 0
// 	limit := 10
// 	total := t

// 	for {
// 		o <- offset
// 		offset = offset + limit
// 		if (offset) > total {
// 			break
// 		}
// 	}
// }

// func (m *Catalog) updateNow(c string) {
// 	m.Updated = time.Now()
// 	mout, err := json.MarshalIndent(m, "", " ")
// 	if err != nil {
// 		logger().Errorln(err)
// 		os.Exit(1)
// 	}
// 	err = ioutil.WriteFile(c, mout, 0644)
// 	if err != nil {
// 		logger().Errorln(err)
// 		os.Exit(1)
// 	}
// 	return

// }

// // UnmarshalJSON is to unf*ck someones bad API.
// func (w *DownloadCountWrapper) UnmarshalJSON(data []byte) (err error) {

// 	if dc, err := strconv.Atoi(string(data)); err == nil {
// 		w.DownloadCount = dc
// 		return nil
// 	}
// 	w.DownloadCount = 0
// 	return nil
// }

// // UnmarshalJSON is to unf*ck someones bad API.
// func (w *InstallCountWrapper) UnmarshalJSON(data []byte) (err error) {

// 	if dc, err := strconv.Atoi(string(data)); err == nil {
// 		w.InstallCount = dc
// 		return nil
// 	}
// 	w.InstallCount = 0
// 	return nil
// }

// func getTotal(s string) (t int) {

// 	var allapps AllApps
// 	res := GetURL(s)
// 	err := json.Unmarshal(res, &allapps)
// 	if err != nil {
// 		logger().Debug(s)
// 		logger().Fatalln(err)
// 		os.Exit(1)
// 	}
// 	return allapps.Total
// }

// // GetURL returns Body from Splunkbase API request.
// func GetURL(u string) (body []byte) {
// 	logger().Debug("Start")
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
// 	}
// 	cl := &http.Client{Transport: tr}
// 	req, err := http.NewRequest("GET", u, nil)
// 	if err != nil {
// 		logger().Fatalln(err)
// 		os.Exit(1)
// 	}
// 	res, err := cl.Do(req)
// 	if err != nil {
// 		logger().Fatalln(err)
// 		os.Exit(1)
// 	}
// 	body, err = ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		logger().Fatalln(err)
// 		os.Exit(1)
// 	}
// 	return body
// }
