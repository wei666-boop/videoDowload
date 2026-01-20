package utils

import (
	"errors"
	"os"
)

//创建随机目录

func RandomID(path string) (string, error) {
	dir, err := os.MkdirTemp(path, "id-*")
	if err != nil {
		return "", errors.New("not create dir")
	}
	return dir, nil
}
