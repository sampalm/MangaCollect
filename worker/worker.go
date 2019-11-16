package worker

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// CompressFiles filename is the output zip file's path.
// files is a list of files to add to the zip.
func CompressFiles(output string, files []string) error {
	zipfile, err := os.OpenFile(output, syscall.O_CREAT|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	for _, file := range files {

		if err := compress(zipWriter, file); err != nil {
			return errors.New("CompressFiles: " + err.Error())
		}
	}
	return nil
}

func compress(zipWriter *zip.Writer, source string) error {

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		fileInfo, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = fmt.Sprintf("%s - %s", filepath.Base(source), info.Name())
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		return err

	})

}

func ReadPath(path, ignore string) []os.FileInfo {
	var templ []os.FileInfo
	list, err := ioutil.ReadDir(path)
	for _, f := range list {
		if !f.IsDir() {
			continue
		}
		if f.Name() == ignore {
			continue
		}
		templ = append(templ, f)
	}
	if err != nil {
		log.Fatal(err)
	}
	return templ
}

func GetLastChapter(path string) (string, error) {
	fl, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer fl.Close()

	if len(fl.File) < 1 {
		return "", err
	}
	cat := strings.Split(fl.File[len(fl.File)-1].Name, "/")
	if len(cat) < 1 {
		return "", fmt.Errorf("GetLastChapter: filepath not found")
	}

	return cat[0], nil
}
