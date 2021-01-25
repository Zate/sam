package sbase

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// TimeTrack is used to measure how long a function takes to run for debug.
func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	pc, file, line, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	fn := name[strings.LastIndex(name, ".")+1:]
	log.WithField("file", filename).WithField("function", fn).WithField("elapsed", fmt.Sprintf("duration: %s", elapsed)).Debugln("Debug")
}

// Logger for logging
func Logger() *log.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}
	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return log.WithField("file", filename).WithField("function", fn)
}

// DoExist checks if the object s exists
func DoExist(s string) bool {
	_, err := os.Stat(s)
	if err != nil {
		return false
	}
	return true
}

// CheckPath looks at the path
func CheckPath(s string) (os.FileInfo, error) {
	// lets validate this path exists
	if _, err := os.Stat(s); err != nil {
		if err := os.MkdirAll(s, 0755); err != nil {
			return nil, err
		}
	}
	dlpathInfo, err := os.Stat(s)
	if err != nil {
		return nil, err
	}
	if dlpathInfo.IsDir() == false {
		return nil, errors.New("Not a directory: " + s)
	}
	return dlpathInfo, nil
}

// ReadFile reads the file at the location of the string given and returns a []byte array of the contents
func ReadFile(f string) ([]byte, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// ChkErr checks for an error.  If there os one, it logs it and fatals
func ChkErr(e error) {
	pc, file, line, _ := runtime.Caller(1)
	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	if e != nil {
		Logger().WithField("file", filename).WithField("function", fn).Errorln(e)
		os.Exit(1)
	}
	return
}
