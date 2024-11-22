package frontend

import _ "embed"

//go:embed dist.zip
var htmlZip []byte

func GetHtmlZipFile() []byte {
	return htmlZip
}
