package tool

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// GeneratedFileName xxxx.pb.gin.go
func GeneratedFileName(oldFilePath string) string {
	dir := filepath.Dir(oldFilePath)
	ext := filepath.Ext(oldFilePath)
	nameWithoutExt := strings.TrimSuffix(filepath.Base(oldFilePath), ext)
	newFileName := fmt.Sprintf("%s.pb.gin.go", nameWithoutExt)

	return filepath.Join(dir, newFileName)
}

// OpenAPI2Go OpenAPI route style convert to go style
func OpenAPI2Go(path string) string {
	if len(path) == 0 {
		return path
	}

	reg := regexp.MustCompile(`^\{(.+)\}$`)
	replacer := func(segment string) string {
		return reg.ReplaceAllString(segment, `:$1`)
	}

	segments := []string{}
	for _, it := range strings.Split(path, "/") {
		segments = append(segments, replacer(it))
	}

	return strings.Join(segments, "/")
}

// ToType 去除包名，留下类型
func ToType(pkgName, inputType string) string {
	return strings.Replace(inputType, "."+pkgName+".", "", 1)
}

// ToServerInterfaceType 服务名称（协议中：{{service}}+GinService）
func ToServerInterfaceType(service *descriptor.ServiceDescriptorProto) string {
	return service.GetName() + "GinServer"
}
