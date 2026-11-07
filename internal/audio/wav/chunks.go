package wav

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// chunks walks the RIFF chunks and returns the format and sample data.
func chunks(r *bytes.Reader) (format, []byte, error) {
	var f format
	var pcm []byte
	var err error
	haveFmt := false
	for {
		var h [8]byte
		if _, err := io.ReadFull(r, h[:]); err != nil {
			break
		}
		id := string(h[0:4])
		size := binary.LittleEndian.Uint32(h[4:8])
		// Streamed WAVs leave the data size at 0 or 0xFFFFFFFF: take
		// whatever is left.
		if rem := uint32(r.Len()); size > rem || (id == "data" && size == 0) {
			size = rem
		}
		body := make([]byte, size)
		if _, err := io.ReadFull(r, body); err != nil {
			return f, nil, err
		}
		if size%2 == 1 {
			r.ReadByte()
		}
		switch id {
		case "fmt ":
			if f, err = parseFmt(body); err != nil {
				return f, nil, err
			}
			haveFmt = true
		case "data":
			pcm = body
		}
	}
	if !haveFmt || pcm == nil {
		return f, nil, errors.New("WAV file has no audio")
	}
	return f, pcm, nil
}

// parseFmt reads a "fmt " chunk.
func parseFmt(body []byte) (format, error) {
	var f format
	if len(body) < 16 {
		return f, errors.New("bad fmt chunk")
	}
	f.tag = binary.LittleEndian.Uint16(body[0:2])
	f.channels = binary.LittleEndian.Uint16(body[2:4])
	f.rate = binary.LittleEndian.Uint32(body[4:8])
	f.bits = binary.LittleEndian.Uint16(body[14:16])
	if f.tag == 0xFFFE && len(body) >= 26 { // WAVE_FORMAT_EXTENSIBLE
		f.tag = binary.LittleEndian.Uint16(body[24:26])
	}
	return f, nil
}

func sample(f format, b []byte) (float32, error) {
	switch {
	case f.tag == 3 && f.bits == 32:
		return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
	case f.tag == 3 && f.bits == 64:
		return float32(math.Float64frombits(binary.LittleEndian.Uint64(b))), nil
	case f.tag == 1 && f.bits == 8:
		return (float32(b[0]) - 128) / 128, nil
	case f.tag == 1 && f.bits == 16:
		return float32(int16(binary.LittleEndian.Uint16(b))) / 32768, nil
	case f.tag == 1 && f.bits == 24:
		return float32(int32(b[0])|int32(b[1])<<8|int32(int8(b[2]))<<16) / 8388608, nil
	case f.tag == 1 && f.bits == 32:
		return float32(int32(binary.LittleEndian.Uint32(b))) / 2147483648, nil
	}
	return 0, fmt.Errorf("unsupported WAV encoding (format %d, %d-bit)", f.tag, f.bits)
}
