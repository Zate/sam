package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Untar takes a dst location and an io.Reader and creates files/directories based on what is in the tar.gz file.
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			logger().Fatalln(err)
			os.Exit(1)
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					logger().Fatalln(err)
					os.Exit(1)
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

	_, err := os.Stat(src)
	if err != nil {
		logger().Fatalln(err)
		os.Exit(1)
	}
	mw := io.MultiWriter(writers...)
	gzw := gzip.NewWriter(mw)
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		f, err := os.Open(file)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		_, err = io.Copy(tw, f)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		f.Close()
		return nil
	})
}

// UnpackApp takes an appid string and finds the latest version of that app and unpacks 3 copies of it.
func UnpackApp(a *App, t []string, f string) {
	logger().Debug("Start")

	logger().Debug(a.Appid)
	tgz := filePath + fmt.Sprint(a.UID) + "/" + a.Appid + "/" + a.LatestVersion + "/" + a.Appid + "_" + a.LatestVersion + ".tar.gz"
	logger().Debug(tgz)

	for _, s := range t {
		target := filePath + fmt.Sprint(a.UID) + "/" + s + "/"
		err := os.MkdirAll(target, 0755)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		r, err := os.Open(tgz)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		err = Untar(target, r)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		var p Package
		p.DType = s
		err = filepath.Walk(target, getFiles(&p))
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		logger().Debug(len(p.Objects))
		a.Packages = append(a.Packages, p)
	}
	if f != "" {
		r, err := os.Open(tgz)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
		err = Untar(f, r)
		if err != nil {
			logger().Fatalln(err)
			os.Exit(1)
		}
	}

	return
}

func getFiles(p *Package) filepath.WalkFunc {
	i := 1
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger().Errorln(err)
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
