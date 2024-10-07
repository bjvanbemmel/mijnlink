package service

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"mime/multipart"
	"os"
)

type FileService struct {
	IndexService IndexService
	FilesDir     string
}

func (s FileService) SaveFile(file multipart.File) (string, error) {
	out, err := os.CreateTemp(s.FilesDir, "*")
	if err != nil {
		return "", err
	}
	defer out.Close()

	buf := bufio.NewWriter(out)
	defer buf.Flush()

	gz := gzip.NewWriter(buf)
	defer gz.Close()
	defer gz.Flush()

	_, err = io.Copy(gz, file)
	if err != nil {
		return "", err
	}

	return s.IndexService.SaveValue(out.Name())
}

func (s FileService) GetFileByKey(key string) (string, error) {
	path, err := s.IndexService.GetValueByKey(key)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}

	buf := bufio.NewReader(gz)
	buffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buffer, buf); err != nil {
		return "", err
	}

	return string(buffer.Bytes()), nil
}
