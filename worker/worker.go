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
		//log.Printf(":: Compressing files %s - %s\n", filepath.Base(source), info.Name())
		_, err = io.Copy(writer, fileToZip)
		return err

	})

}

func ReadPath(path string) []os.FileInfo {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return list
}
