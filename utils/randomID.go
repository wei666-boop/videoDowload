package utils

import (
	"errors"
	"os"
)

//创建随机目录

func RandomID() (string, error) {
	dir, err := os.MkdirTemp("temp", "id-*")
	if err != nil {
		return "", errors.New("not create dir")
	}
	return dir, nil
}
