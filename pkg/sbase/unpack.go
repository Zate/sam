package sbase

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Untar takes a dst location and an io.Reader and creates files/directories based on what is in the tar.gz file.
func Untar(dst string, r io.Reader) error {
	defer TimeTrack(time.Now())
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}

// Tar takes a src folder as a string and writes out a .tar.gz file
func Tar(src string, writers ...io.Writer) error {
	defer TimeTrack(time.Now())
	_, err := os.Stat(src)
	if err != nil {
		Logger().Fatalln(err)
		os.Exit(1)
	}
	mw := io.MultiWriter(writers...)
	gzw := gzip.NewWriter(mw)
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		f, err := os.Open(file)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		_, err = io.Copy(tw, f)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		f.Close()
		return nil
	})
}

// UnpackApp takes an appid string and finds the latest version of that app and unpacks 3 copies of it.
func UnpackApp(a *App, t []string, f string) {
	defer TimeTrack(time.Now())

	tgz := FilePath + fmt.Sprint(a.UID) + "/" + a.Appid + "/" + a.LatestVersion + "/" + a.Appid + "_" + a.LatestVersion + ".tar.gz"

	for _, s := range t {
		target := FilePath + fmt.Sprint(a.UID) + "/" + s + "/"
		err := os.MkdirAll(target, 0755)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		r, err := os.Open(tgz)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		err = Untar(target, r)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		var p Package
		p.DType = s
		err = filepath.Walk(target, GetFiles(&p))
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		a.Packages = append(a.Packages, p)
	}
	if f != "" {
		r, err := os.Open(tgz)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
		err = Untar(f, r)
		if err != nil {
			Logger().Fatalln(err)
			os.Exit(1)
		}
	}
	return
}

// GetFiles walks all the folders in a package and puts the file info into the Package struct
func GetFiles(p *Package) filepath.WalkFunc {
	defer TimeTrack(time.Now())
	i := 1
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Logger().Errorln(err)
			return err
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
