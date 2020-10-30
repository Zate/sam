package pkgr

import (
	"fmt"
	"os"
	"path/filepath"
)

func parsePackages(t []string, a *App) error {

	for _, s := range t {
		var p Package
		p.DType = s
		err := filepath.Walk(filePath+fmt.Sprint(a.UID)+"/"+s+"/", getFiles(&p))
		if err != nil {
			return err
		}
		logger().Debug(len(p.Objects))
		a.Packages = append(a.Packages, p)
	}

	return nil
}

func getFiles(p *Package) filepath.WalkFunc {
	i := 1
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger().Errorln(err)
			return nil
		}
		var o FSObject
		o.Type = "file"
		if info.Mode().IsRegular() != true {
			o.Type = "weird"
		}
		o.ID = i
		o.Name = info.Name()
		o.RelativePath = path
		o.FileInfo = info
		if info.IsDir() {
			o.Type = "dir"
		}
		p.Objects = append(p.Objects, o)
		i++
		return nil
	}
}
