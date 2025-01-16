package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipData(Str []byte) (err error, res []byte) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(Str)
	gz.Close()
	return nil, buf.Bytes()
}

func UnGzipData(Str []byte) (RStr []byte) {
	rdata := bytes.NewReader(Str)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	return s
}