package util

import (
	"archive/zip"
	"bytes"
)

func Zip(filelist map[string][]byte) (zipdata []byte, err error) {
	buff := new(bytes.Buffer)
	w := zip.NewWriter(buff)
	for filename, data := range filelist {
		f, err := w.Create(filename)
		if err != nil {
			return zipdata, err
		}
		_, err = f.Write([]byte(data))
		if err != nil {
			return zipdata, err
		}
	}
	err = w.Close()
	zipdata = buff.Bytes()
	return
}
