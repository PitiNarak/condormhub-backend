package utils

import "os"

func ReadTemplate(path string) (string, error) {
	byteContent, err := os.ReadFile(path)
	if err != nil { //many people wrap this into a function
		return "", err
	}
	return string(byteContent), nil
}
