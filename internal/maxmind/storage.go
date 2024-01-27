package maxmind

import (
	"errors"
	"os"
	"path/filepath"
)

const fileName = "GeoLite2-City.mmdb"

func FindDBLocation() (string, error) {
	var filePath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == fileName {
			filePath = path
		}
		return nil
	})
	if filePath == "" {
		return "", errors.New("cannot find mmdb file")
	}
	return filePath, err
}
