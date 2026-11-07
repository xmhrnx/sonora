package studio

import (
	"strings"

	"github.com/mehran9675/sonora/internal/textfile"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/tr"
)

// openFile replaces the script with the text of a document the user
// picks: text, Markdown, subtitles, HTML, Word, OpenDocument, EPUB, PDF
// or RTF.
func (e *editor) openFile(c *core.Ctx) {
	var text, name string
	c.Background(func() error {
		f, err := c.PickFile(textfile.Extensions, "")
		if err != nil {
			return err
		}
		defer f.Close()
		name = f.Name
		text, err = textfile.Read(name, f)
		return err
	}, func(err error) {
		if err != nil {
			c.Error(err)
			return
		}
		e.in.ChangeText(text)
		words := len(strings.Fields(text))
		msg := tr.N("Loaded {n} words", words)
		if name != "" {
			msg = tr.N("Loaded {n} words from {file}", words, "file", name)
		}
		c.Success(msg)
	})
}
