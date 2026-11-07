package textfile

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// zipKind tells DOCX, ODT and EPUB apart by their contents.
func zipKind(data []byte) string {
	z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "txt"
	}
	for _, f := range z.File {
		switch f.Name {
		case "word/document.xml":
			return "docx"
		case "content.xml":
			return "odt"
		case "META-INF/container.xml":
			return "epub"
		}
	}
	return "zip"
}

func openZip(data []byte) (*zip.Reader, error) {
	z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, errors.New("the document is damaged or not what its name says")
	}
	return z, nil
}

func zipFile(z *zip.Reader, name string) ([]byte, error) {
	f, err := z.Open(name)
	if err != nil {
		return nil, fmt.Errorf("the document has no %s", name)
	}
	defer f.Close()
	return io.ReadAll(io.LimitReader(f, maxSize))
}
