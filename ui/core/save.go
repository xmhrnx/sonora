package core

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/tr"
)

// saveAs asks where to save with the system's dialog; a desktop
// without one saves into Downloads.
func (c *Ctx) saveAs(name string, open func() (io.ReadCloser, error)) {
	var fallback string
	c.Background(func() error {
		w, err := c.App.SaveFile(context.Background(), name)
		if err != nil && !errors.Is(err, giode.ErrCanceled) && !Mobile {
			fallback = filepath.Join(downloadsDir(), name)
			w, err = os.Create(fallback)
		}
		if err != nil {
			return err
		}
		in, err := open()
		if err != nil {
			w.Close()
			return err
		}
		defer in.Close()
		if _, err := io.Copy(w, in); err != nil {
			w.Close()
			return err
		}
		return w.Close()
	}, func(err error) {
		switch {
		case err != nil:
			c.Error(err)
		case fallback != "":
			c.Success(tr.T("Saved to {path}", "path", fallback))
		default:
			c.Success(tr.T("Saved"))
		}
	})
}
