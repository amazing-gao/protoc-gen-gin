package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func FillPath(oldFilePath string) string {
	dir := filepath.Dir(oldFilePath)
	ext := filepath.Ext(oldFilePath)
	nameWithoutExt := strings.TrimSuffix(filepath.Base(oldFilePath), ext)
	newFileName := fmt.Sprintf("%s.pb.gin.go", nameWithoutExt)

	return filepath.Join(dir, newFileName)
}

func totype(pkgName, inputType string) string {
	return strings.Replace(inputType, "."+pkgName+".", "", 1)
}
