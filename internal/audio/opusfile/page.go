package opusfile

import "encoding/binary"

// addPacket queues a packet for the current page.
func (w *Writer) addPacket(p []byte) {
	for n := len(p); ; n -= 255 {
		if n < 255 {
			w.lacing = append(w.lacing, byte(n))
			break
		}
		w.lacing = append(w.lacing, 255)
	}
	w.pageData = append(w.pageData, p...)
	w.packets++
}

func (w *Writer) flushPage(flags byte) error {
	if w.packets == 0 && flags == 0 {
		return nil
	}
	err := w.writeRaw(w.lacing, w.pageData, w.granule, flags)
	w.lacing, w.pageData, w.packets = w.lacing[:0], w.pageData[:0], 0
	if err != nil {
		return err
	}
	// Push each page to the file, so a crash loses at most a second.
	return w.w.Flush()
}

// writePage writes a page holding whole packets.
func (w *Writer) writePage(packets [][]byte, granule uint64, flags byte) error {
	var lacing, data []byte
	for _, p := range packets {
		for n := len(p); ; n -= 255 {
			if n < 255 {
				lacing = append(lacing, byte(n))
				break
			}
			lacing = append(lacing, 255)
		}
		data = append(data, p...)
	}
	return w.writeRaw(lacing, data, granule, flags)
}

func (w *Writer) writeRaw(lacing, data []byte, granule uint64, flags byte) error {
	h := make([]byte, 27, 27+len(lacing))
	copy(h, "OggS")
	h[5] = flags
	binary.LittleEndian.PutUint64(h[6:], granule)
	binary.LittleEndian.PutUint32(h[14:], w.serial)
	binary.LittleEndian.PutUint32(h[18:], w.seq)
	h[26] = byte(len(lacing))
	h = append(h, lacing...)
	crc := crcUpdate(crcUpdate(0, h), data)
	binary.LittleEndian.PutUint32(h[22:], crc)
	w.seq++
	if _, err := w.w.Write(h); err != nil {
		return err
	}
	_, err := w.w.Write(data)
	return err
}
