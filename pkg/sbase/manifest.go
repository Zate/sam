package sbase

import (
	"encoding/json"
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
		sb.Manifest.updateNow()
		return
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

func (m *Manifest) updateNow() {
	m.Updated = time.Now()
	mout, err := json.MarshalIndent(m, "", " ")
	ChkErr(err)
	err = ioutil.WriteFile(FilePath+"/manifest.json", mout, 0644)
	ChkErr(err)
	return
}
