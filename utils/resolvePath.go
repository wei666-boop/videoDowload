package utils

import (
	"path/filepath"
)

//判断是否是绝对路径，如果是的话就使用绝对路径，否则使用相对路径+appPath

func ResolvePath(appPath string, p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(appPath, p)
	}
	return p
}
