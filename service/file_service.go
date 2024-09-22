package service

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type FileService struct{}

func (s FileService) SaveFile(file multipart.File) (string, error) {
	out, err := os.CreateTemp(".files/", "*")
	if err != nil {
		return "", err
	}

	buf := bufio.NewWriter(out)
	defer buf.Flush()

	gz := gzip.NewWriter(buf)
	defer gz.Flush()

	i, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d bytes written", i), nil
}
