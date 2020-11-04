package main

// import (
// 	"os"
// 	"runtime"
// 	"strconv"
// 	"strings"

// 	log "github.com/sirupsen/logrus"
// )

// func logger() *log.Entry {
// 	pc, file, line, ok := runtime.Caller(1)
// 	if !ok {
// 		panic("Could not get context info for logger!")
// 	}

// 	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
// 	funcname := runtime.FuncForPC(pc).Name()
// 	fn := funcname[strings.LastIndex(funcname, ".")+1:]
// 	return log.WithField("file", filename).WithField("function", fn)
// }

// // doExist checks if the object s exists
// func doExist(s string) bool {
// 	_, err := os.Stat(s)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }
