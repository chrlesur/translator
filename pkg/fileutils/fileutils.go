package fileutils

import (
	"io/ioutil"
	"os"
)

func ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
