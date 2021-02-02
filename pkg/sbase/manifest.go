package sbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// LoadManifest checks if ther is an apps/manifest.json and if there is, parses it into a new Manifest
// If there is not, it creates it, and starts a new empty Manifest
func (sb *SBase) LoadManifest(m string) {
	defer TimeTrack(time.Now())
	if DoExist(m) == false {
		if err := os.MkdirAll(m, 0755); err != nil {
			ChkErr(err)
		}
	}
	if DoExist(m+"/manifest.json") == false {
		// fmt.Sprint("t")
		f, err := os.OpenFile(m+"/manifest.json", os.O_RDONLY|os.O_CREATE, 0644)
		ChkErr(err)
		err = f.Close()
		ChkErr(err)
		sb.UpdateNow(m)
		return
	} else {

	}
	f, err := ioutil.ReadFile(m + "/manifest.json")
	ChkErr(err)
	err = json.Unmarshal([]byte(f), &sb.Manifest)
	ChkErr(err)
	return
}

func (m *Manifest) setNow() {
	m.Updated = time.Now()
}

// UpdateNow writes sb.Manifest out to file
func (sb *SBase) UpdateNow(m string) {
	sb.Manifest.Updated = time.Now()
	mout, err := json.MarshalIndent(&sb.Manifest, "", " ")
	ChkErr(err)
	fmt.Print(m + "/manifest.json")
	err = ioutil.WriteFile(m+"/manifest.json", mout, 0755)
	ChkErr(err)
	return
}

// CheckApps compares what is in manifest.json to what is on disk in FilePath
func (sb *SBase) CheckApps(fp string) {
	// Read app uid's from manifest.json

	// Read app uid's from FilePath

	// Validate appinfo.json against the struct and check if it's in the manifest.json uids

	// If it's in the manifest but not on disk, download it.

	// if its on disk but not in the manifest, and the appinfo.json is valid, check if we have the latest version.

	// if not latest version, downloadlatest version

	return
}
