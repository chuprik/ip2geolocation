package maxmind

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

const downloadUrl = "https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=%s&suffix=tar.gz"

func Download(editionId, licenseKey string) error {
	url := fmt.Sprintf(downloadUrl, editionId, licenseKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return Unpack(resp.Body)
}

func Unpack(stream io.ReadCloser) error {
	gzr, err := gzip.NewReader(stream)
	if err != nil {
		return nil
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.Mkdir(header.Name, 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				return err

			}
			outFile.Close()
		}
	}
}
