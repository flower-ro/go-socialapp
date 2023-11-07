package utils

import (
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func CopyFile(sourceFile string, dst string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, input, 0644)
}
func CopyFileByOs(sourceFile string, dst string) error {
	source, err := os.Open(sourceFile) //open the source file
	if err != nil {
		panic(err)
	}
	defer source.Close()

	destination, err := os.Create(dst) //create the destination file
	if err != nil {
		panic(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source) //copy the contents of source to destination file
	return err
}

// RemoveFile is removing file with delay
func RemoveFile(delaySecond int, paths ...string) error {
	if delaySecond > 0 {
		time.Sleep(time.Duration(delaySecond) * time.Second)
	}

	for _, path := range paths {
		if path != "" {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			//目录不存在
			return errors.WithCode(code.DirNotExistedErr, err.Error())
		} else {
			return errors.WithCode(code.NotSureIsExisted, err.Error())
		}

	}
	return nil
}

func CreateFile(dir string, fileName string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			//目录不存在
			return errors.WithCode(code.DirNotExistedErr, err.Error())
		} else {
			return errors.WithCode(code.NotSureIsExisted, err.Error())
		}

	}
	newFolder := filepath.Join(dir, fileName)
	_, err = os.Stat(newFolder)
	if err == nil {
		return errors.WithCode(code.FileIsExisted, err.Error())
	}

	_, err = os.Create(newFolder)
	if err != nil {
		return errors.WithCode(code.FileCreatedFail, err.Error())
	}
	return nil
}

// CreateFolder create new folder and sub folder if not exist
func CreateFolder(folderPath ...string) error {
	for _, folder := range folderPath {
		newFolder := filepath.Join(".", folder)
		err := os.MkdirAll(newFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
