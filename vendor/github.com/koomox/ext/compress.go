package ext

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// Decompress .tar.gz file
func DeCompress(tarFile, dest string) ([]string, error) {
	completeFile := make([]string, 0)
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return completeFile, err
	}
	defer srcFile.Close()

	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return completeFile, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return completeFile, err
			}
		}
		filename := dest + hdr.Name
		if strings.HasSuffix(filename, "/") { // Dir
			err := createDir(filename)
			if err != nil {
				return completeFile, err
			}
		} else { // file
			file, err := createFile(filename)
			if err != nil {
				return completeFile, err
			}
			io.Copy(file, tr)
		}
		completeFile = append(completeFile, filename)
	}
	return completeFile, nil
}

func createDir(name string) (err error) {
	path := string([]rune(name)[0:strings.LastIndex(name, "/")])
	exist, err := PathExist(path)
	if err != nil {
		return
	}
	if !exist {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	if err := createDir(name); err != nil {
		return nil, err
	}
	return os.Create(name)
}
