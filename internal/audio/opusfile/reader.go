package opusfile

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/thesyncim/gopus"
	"github.com/thesyncim/gopus/container/ogg"
)

// ErrNotOpus means the file isn't Ogg Opus.
var ErrNotOpus = errors.New("not an Ogg Opus file")

// IsOpus reports whether head, the start of a file, is Ogg Opus.
func IsOpus(head []byte) bool {
	return len(head) >= 36 && string(head[:4]) == "OggS" && string(head[28:36]) == "OpusHead"
}

// Decode reads an Ogg Opus stream, calling fn with mono 48 kHz samples
// a packet at a time.
func Decode(r io.Reader, fn func(rate int, mono []float32) error) error {
	or, err := ogg.NewReader(r)
	if err != nil {
		return ErrNotOpus
	}
	ch := max(1, int(or.Channels()))
	if ch > 2 {
		return errors.New("Opus files with more than two channels aren't supported")
	}
	dec, err := gopus.NewDecoder(gopus.DefaultDecoderConfig(48000, ch))
	if err != nil {
		return err
	}
	skip := int(or.PreSkip())
	pcm := make([]float32, 5760*ch) // the longest Opus packet: 120 ms
	mono := make([]float32, 5760)
	for {
		p, _, err := or.ReadPacket()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n, err := dec.Decode(p, pcm)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			var sum float32
			for c := 0; c < ch; c++ {
				sum += pcm[i*ch+c]
			}
			mono[i] = sum / float32(ch)
		}
		out := mono[:n]
		if skip > 0 {
			d := min(skip, len(out))
			out, skip = out[d:], skip-d
		}
		if len(out) > 0 {
			if err := fn(48000, out); err != nil {
				return err
			}
		}
	}
}

// opusHead builds the identification header (RFC 7845 §5.1) for mono.
func opusHead(preSkip, rate int) []byte {
	b := make([]byte, 19)
	copy(b, "OpusHead")
	b[8] = 1 // version
	b[9] = 1 // channels
	binary.LittleEndian.PutUint16(b[10:], uint16(preSkip))
	binary.LittleEndian.PutUint32(b[12:], uint32(rate))
	return b // output gain 0, mapping family 0
}

// opusTags builds the comment header (RFC 7845 §5.2) with no comments.
func opusTags(vendor string) []byte {
	b := make([]byte, 8+4+len(vendor)+4)
	copy(b, "OpusTags")
	binary.LittleEndian.PutUint32(b[8:], uint32(len(vendor)))
	copy(b[12:], vendor)
	return b
}
