package datauri

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func ReadDataURI(contentType string, file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	encodedData := base64.StdEncoding.EncodeToString(bs)

	datauri := fmt.Sprintf("data:%s;base64,%s", contentType, encodedData)

	return datauri, nil
}
