package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func CheckAndCreatePath(path string) error {
	if _, err := os.Stat(path); err != nil {
		if !os.IsExist(err) { //目录不存在，创建目录
			return os.MkdirAll(path, os.ModePerm)
		}
	}
	return nil
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
