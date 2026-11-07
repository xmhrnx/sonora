package decode

import (
	"io"

	"github.com/mewkiz/flac"
)

func decodeFLAC(r io.Reader, out *resampler) error {
	s, err := flac.New(r)
	if err != nil {
		return err
	}
	rate := int(s.Info.SampleRate)
	scale := float32(int64(1) << (s.Info.BitsPerSample - 1))
	var mono []float32
	for {
		fr, err := s.ParseNext()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n := fr.Subframes[0].NSamples
		if cap(mono) < n {
			mono = make([]float32, n)
		}
		mono = mono[:n]
		clear(mono)
		for _, sub := range fr.Subframes {
			for i, v := range sub.Samples[:n] {
				mono[i] += float32(v)
			}
		}
		k := scale * float32(len(fr.Subframes))
		for i := range mono {
			mono[i] /= k
		}
		out.write(rate, mono)
	}
}
