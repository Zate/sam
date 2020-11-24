package main

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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

// doExist checks if the object s exists
func doExist(s string) bool {
	_, err := os.Stat(s)
	if err != nil {
		return false
	}
	return true
}

func checkPath(s string) (os.FileInfo, error) {
	// lets validate this path exists
	dlpathInfo, err := os.Stat(s)
	if err != nil {
		return nil, err
	}
	if dlpathInfo.IsDir() == false {
		return nil, errors.New("Not a directory: " + s)
	}
	return dlpathInfo, nil
}
