package onnx

import (
	"fmt"
	"io"

	"github.com/mehran9675/sonora/internal/engines"
)

// progressReader logs "label NN%" as bytes arrive, in 2% steps.
type progressReader struct {
	r     io.Reader
	total int64
	read  int64
	step  int64
	label string
	log   engines.Logger
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	if p.total > 0 {
		if s := p.read * 50 / p.total; s != p.step {
			p.step = s
			p.log(fmt.Sprintf("%s… %d%%", p.label, p.read*100/p.total))
		}
	}
	return n, err
}
