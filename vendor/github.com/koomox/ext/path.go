package ext

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var rawFileNameUniqueErr = errors.New("Generator Unique File Name Failed!")

const (
	rawFileNameLength = 32
)

// is Dir And file Exists Return true
func PathExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetCurrentDirectory() (string, error) {
	return currentDirectory()
}

// Get Now Dir
func currentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", fmt.Errorf("Get Current Directory Err:%v", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func GetCustomDirectoryAllFile(path string) ([]string, error) {
	var completeList []string
	return customDirectoryAllFile(path, completeList)
}

// @ path string
// @ completeList  while param, var completeList []string
func customDirectoryAllFile(path string, completeList []string) ([]string, error) {
	var err error
	if path == "" {
		if path, err = currentDirectory(); err != nil {
			return nil, err
		}
	}
	List, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range List {
		if f.IsDir() {
			childPath := path + "/" + f.Name()
			completeList, err = customDirectoryAllFile(childPath, completeList)
			if err != nil {
				return nil, err
			}
		} else {
			pathName := path + "/" + f.Name()
			completeList = append(completeList, pathName)
		}
	}
	return completeList, nil
}

func findCustomFile(fList []string, suffix string) (string, bool) {
	suffixName := strings.ToLower(suffix)
	for _, f := range fList {
		if strings.HasSuffix(strings.ToLower(f), suffixName) {
			return f, true
		}
	}
	return "", false
}

func FindCurrentDirectoryFile(suffix string) (string, bool) {
	path, err := currentDirectory()
	if err != nil {
		return "", false
	}
	var fList []string
	fList, err = customDirectoryAllFile(path, fList)
	if err != nil {
		return "", false
	}

	return findCustomFile(fList, suffix)
}

// 参数为空，检查当前目录
func GeneratorRawFileNameUnique(path string) (string, error) {
	var err error
	if path == "" {
		path, err = currentDirectory()
		if err != nil {
			return "", err
		}
	}
	for i := 0; i < 3; i++ {
		f, err := RandomString(rawFileNameLength)
		if err != nil {
			return "", err
		}
		fn := path + "/" + f
		exist, err := PathExist(fn)
		if err != nil || exist {
			continue
		}
		return f, nil
	}

	return "", rawFileNameUniqueErr
}