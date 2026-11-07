package decode

import (
	"io"

	"github.com/jfreymuth/oggvorbis"
)

func decodeVorbis(r io.Reader, out *resampler) error {
	d, err := oggvorbis.NewReader(r)
	if err != nil {
		return err
	}
	ch := d.Channels()
	buf := make([]float32, ch*4096)
	mono := make([]float32, 4096)
	for {
		n, err := d.Read(buf)
		frames := n / ch
		for i := 0; i < frames; i++ {
			var sum float32
			for c := 0; c < ch; c++ {
				sum += buf[i*ch+c]
			}
			mono[i] = sum / float32(ch)
		}
		if frames > 0 {
			out.write(d.SampleRate(), mono[:frames])
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
