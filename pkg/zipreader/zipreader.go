package zipreader

import (
	"archive/zip"
	"strings"
)

func ReadZip(fileName string, ext string) ([]string, error) {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res := make([]string, 0)
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ext) {
			res = append(res, f.Name)
		}
	}
	return res, nil
}
