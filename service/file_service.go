package service

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type FileService struct {
	IndexService IndexService
	FilesDir     string
}

type FileMetaData struct {
	FileName string `json:"filename,omitempty"`
}

func (s FileService) SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
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

	return s.IndexService.SaveValue(out.Name(), header)
}

func (s FileService) GetFileByKey(key string) (string, string, error) {
	value, err := s.IndexService.GetValueByKey(key)
	if err != nil {
		return "", "", err
	}

	split := strings.Split(value, META_DELIMITER)
	path := split[0]
	buffer, err := s.getBufferFromFileByPath(path)

	if len(split) < 2 {
		return buffer.String(), fmt.Sprintf("%d", time.Now().Unix()), err
	}

	metaRaw := split[1]
	var meta FileMetaData
	err = json.Unmarshal([]byte(metaRaw), &meta)
	if err != nil {
		return buffer.String(), "", err
	}

	return buffer.String(), meta.FileName, err
}

func (s FileService) getBufferFromFileByPath(path string) (*bytes.Buffer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(gz)
	buffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buffer, buf); err != nil {
		return nil, err
	}

	return buffer, nil
}
