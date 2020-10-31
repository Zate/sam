package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Manifest is a struct representing a json file with a list of all the apps installed and managed
type Manifest struct {
	Apps    []App     `json:"apps"`
	Updated time.Time `json:"updated"`
}

// LoadManifest checks if ther is an apps/manifest.json and if there is, parses it into a new Manifest
// If there is not, it creates it, and starts a new empty Manifest
func LoadManifest(m string) {

	M := Manifest{
		[]App{},
		time.Now(),
	}
	if doExist(filePath+"manifest.json") == false {
		_, err := os.Create(filePath + "manifest.json")
		if err != nil {
			logger().Errorln(err)
			os.Exit(1)
		}
		M.updateNow()
		return
	}
	f, err := ioutil.ReadFile(m)
	if err != nil {
		logger().Errorln(err)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(f), &M)
	if err != nil {
		logger().Errorln(err)
	}
	return
}

func (m *Manifest) setNow() {
	m.Updated = time.Now()
}

func (m *Manifest) updateNow() {
	m.Updated = time.Now()
	mout, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		logger().Errorln(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filePath+"/manifest.json", mout, 0644)
	if err != nil {
		logger().Errorln(err)
		os.Exit(1)
	}
	return

}
