package textfile

import (
	"bytes"
	"encoding/xml"
	"errors"
	"path"
	"strings"
)

// xmlText walks an XML document: text from character data inside
// elements named textEl, a line break at the end of elements in
// breakEl, and replacements for empty elements in emptyEl.
func xmlText(data []byte, textEl string, breakEl map[string]bool, emptyEl map[string]string) string {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	var b strings.Builder
	depth := 0 // inside textEl (any text when textEl is "")
	for {
		tok, err := d.Token()
		if err != nil {
			return b.String()
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == textEl {
				depth++
			}
			if r, ok := emptyEl[t.Name.Local]; ok {
				b.WriteString(r)
			}
		case xml.EndElement:
			if t.Name.Local == textEl && depth > 0 {
				depth--
			}
			if breakEl[t.Name.Local] {
				b.WriteString("\n")
			}
		case xml.CharData:
			if textEl == "" || depth > 0 {
				b.Write(t)
			}
		}
	}
}

func docx(data []byte) (string, error) {
	z, err := openZip(data)
	if err != nil {
		return "", err
	}
	doc, err := zipFile(z, "word/document.xml")
	if err != nil {
		return "", err
	}
	return xmlText(doc, "t", map[string]bool{"p": true}, map[string]string{"tab": " ", "br": "\n", "cr": "\n"}), nil
}

func odt(data []byte) (string, error) {
	z, err := openZip(data)
	if err != nil {
		return "", err
	}
	doc, err := zipFile(z, "content.xml")
	if err != nil {
		return "", err
	}
	// Only the body: styles and fonts sit in office:automatic-styles,
	// which holds no character data, so all text is the document's.
	return xmlText(doc, "", map[string]bool{"p": true, "h": true}, map[string]string{"s": " ", "tab": " ", "line-break": "\n"}), nil
}

// epub reads the book's chapters in reading order.
func epub(data []byte) (string, error) {
	z, err := openZip(data)
	if err != nil {
		return "", err
	}
	container, err := zipFile(z, "META-INF/container.xml")
	if err != nil {
		return "", err
	}
	var c struct {
		Rootfiles []struct {
			Path string `xml:"full-path,attr"`
		} `xml:"rootfiles>rootfile"`
	}
	if xml.Unmarshal(container, &c) != nil || len(c.Rootfiles) == 0 {
		return "", errors.New("the book has no table of contents")
	}
	opfPath := c.Rootfiles[0].Path
	opfData, err := zipFile(z, opfPath)
	if err != nil {
		return "", err
	}
	var opf struct {
		Items []struct {
			ID   string `xml:"id,attr"`
			Href string `xml:"href,attr"`
		} `xml:"manifest>item"`
		Spine []struct {
			IDRef string `xml:"idref,attr"`
		} `xml:"spine>itemref"`
	}
	if err := xml.Unmarshal(opfData, &opf); err != nil {
		return "", errors.New("the book's contents can't be read")
	}
	hrefs := map[string]string{}
	for _, it := range opf.Items {
		hrefs[it.ID] = it.Href
	}
	var parts []string
	for _, ref := range opf.Spine {
		href, ok := hrefs[ref.IDRef]
		if !ok {
			continue
		}
		name := path.Join(path.Dir(opfPath), href)
		if i := strings.IndexByte(name, '#'); i >= 0 {
			name = name[:i]
		}
		chapter, err := zipFile(z, name)
		if err != nil {
			continue
		}
		parts = append(parts, html(chapter))
	}
	return strings.Join(parts, "\n\n"), nil
}
