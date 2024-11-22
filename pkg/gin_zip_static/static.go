package gin_zip_static

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

var fileContentTypeMap = []struct {
	Ext         string
	ContentType string
}{
	{
		Ext:         ".js",
		ContentType: "application/javascript",
	},
	{
		Ext:         ".css",
		ContentType: "text/css",
	},
	{
		Ext:         "manifest",
		ContentType: "application/octet-stream",
	},
	{
		Ext:         ".png",
		ContentType: "image/png",
	},
	{
		Ext:         ".jpg",
		ContentType: "image/jpeg",
	},
	{
		Ext:         ".jpeg",
		ContentType: "image/jpeg",
	},
	{
		Ext:         "",
		ContentType: "text/html; charset=utf-8",
	},
}

// MatchFile 匹配不同文件类型的Content-Type
func MatchFile(fileName string) (result string) {
	result = "text/html; charset=utf-8"
	for _, value := range fileContentTypeMap {
		if strings.Contains(fileName, value.Ext) {
			result = value.ContentType
			break
		}
	}
	return result
}

func RegisterStaticFile(zipFile []byte, callback func(fileMap map[string][]byte)) {
	var result = make(map[string][]byte)
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		panic("加载静态资源失败" + err.Error())
	}
	for _, file := range zipReader.File {
		open, err := file.Open()
		if err != nil {
			panic("读取静态资源失败" + err.Error())
		}
		data, err := io.ReadAll(open)
		if err != nil {
			panic("读取静态资源失败" + err.Error())
		}
		_ = open.Close()
		result[file.Name] = data
	}
	callback(result)
}
