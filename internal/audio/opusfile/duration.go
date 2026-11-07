package opusfile

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"time"
)

// Duration reads an Ogg Opus file's length from its header and last
// page, without decoding it.
func Duration(path string) (time.Duration, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	head := make([]byte, 64)
	if n, _ := io.ReadFull(f, head); !IsOpus(head[:n]) {
		return 0, ErrNotOpus
	}
	preSkip := int64(binary.LittleEndian.Uint16(head[28+10:]))
	st, err := f.Stat()
	if err != nil {
		return 0, err
	}
	tail := min(st.Size(), 64<<10)
	buf := make([]byte, tail)
	if _, err := f.ReadAt(buf, st.Size()-tail); err != nil && err != io.EOF {
		return 0, err
	}
	for i := len(buf) - 27; i >= 0; i-- {
		if string(buf[i:i+4]) == "OggS" && buf[i+4] == 0 {
			granule := int64(binary.LittleEndian.Uint64(buf[i+6:]))
			if granule < 0 { // -1: no packet ends on this page
				continue
			}
			return time.Duration(max(0, granule-preSkip)) * time.Second / 48000, nil
		}
	}
	return 0, errors.New("no Ogg pages found")
}
