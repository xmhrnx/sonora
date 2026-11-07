package decode

import (
	"io"

	"github.com/hajimehoshi/go-mp3"
)

func decodeMP3(r io.Reader, out *resampler) error {
	d, err := mp3.NewDecoder(r)
	if err != nil {
		return err
	}
	// go-mp3 always produces 16-bit little-endian stereo.
	buf := make([]byte, 4*4096)
	mono := make([]float32, 4096)
	for {
		n, err := io.ReadFull(d, buf)
		frames := n / 4
		for i := 0; i < frames; i++ {
			l := int16(uint16(buf[4*i]) | uint16(buf[4*i+1])<<8)
			r := int16(uint16(buf[4*i+2]) | uint16(buf[4*i+3])<<8)
			mono[i] = (float32(l) + float32(r)) / 65536
		}
		if frames > 0 {
			out.write(d.SampleRate(), mono[:frames])
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
