package sandbox

import (
	"archive/tar"
	"bytes"
	"io"
)

// CreateTarArchive is a function that Generates a valid tarball entirely in memory using a bytes.Buffer.
func CreateTarArchive(fileName, fileContent string) (io.Reader, error) {
	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{Name: fileName, Size: int64(len(fileContent)), Mode: 0644}

	tw.WriteHeader(hdr)
	tw.Write([]byte(fileContent))
	tw.Close()

	return &buf, nil
}
