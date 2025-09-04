package packer

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pm/config"
)

// CollectFiles gathers files based on target patterns, excluding specified patterns
func CollectFiles(targets []interface{}) ([]string, error) {
	var files []string
	for _, t := range targets {
		target := t.(config.Target)
		matches, err := filepath.Glob(target.Path)
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			if target.Exclude != "" {
				excluded, _ := filepath.Match(target.Exclude, filepath.Base(match))
				if excluded {
					continue
				}
			}
			files = append(files, match)
		}
	}
	return files, nil
}

// CreateZip creates a ZIP archive from files, returns buffer
func CreateZip(files []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	for _, file := range files {
		relPath, err := filepath.Rel(".", file)
		if err != nil {
			return nil, err
		}
		zw, err := zipWriter.Create(relPath)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		_, err = zw.Write(content)
		if err != nil {
			return nil, err
		}
	}
	zipWriter.Close()
	return buf, nil
}

// Unzip extracts a ZIP file to the destination directory
func Unzip(zipFile, dest string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		err = outFile.Close()
		if err != nil {
			return err
		}
		err = rc.Close()
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}
	return nil
}
