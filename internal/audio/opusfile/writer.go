// Package opusfile writes and reads Ogg Opus files (RFC 7845), the
// format Sonora records in: speech at 16-64 kbps takes a few percent
// of a WAV's size. The codec is gopus, a pure-Go Opus implementation
// checked against libopus, so no C library is needed on any platform.
package opusfile

import (
	"bufio"
	"math/rand/v2"
	"os"

	"github.com/thesyncim/gopus"
)

// Writer encodes mono audio into an Ogg Opus file as it arrives.
type Writer struct {
	f    *os.File
	w    *bufio.Writer
	enc  *gopus.Encoder
	rate int // input sample rate

	frame   []float32 // samples waiting for a full frame
	packet  []byte
	preSkip int

	serial  uint32
	seq     uint32
	granule uint64 // 48 kHz samples encoded so far, pre-skip included
	written uint64 // input samples written, for the final granule

	pageData []byte
	lacing   []byte
	packets  int
}

// frameMS is the length of each Opus frame.
const frameMS = 20

// packetsPerPage groups about a second of audio into each Ogg page, so
// the container costs little and an interrupted recording loses at most
// that much.
const packetsPerPage = 50

// Create starts a mono Ogg Opus file at path. rate is the input sample
// rate (8000, 12000, 16000, 24000 or 48000) and bitrate the target in
// bits per second.
func Create(path string, rate, bitrate int) (*Writer, error) {
	// The audio application keeps the full band, the better choice for
	// voices that will be cloned.
	enc, err := gopus.NewEncoder(gopus.EncoderConfig{SampleRate: rate, Channels: 1, Application: gopus.ApplicationAudio})
	if err != nil {
		return nil, err
	}
	if err := enc.SetBitrate(bitrate); err != nil {
		return nil, err
	}
	if err := enc.SetFrameSize(rate * frameMS / 1000); err != nil {
		return nil, err
	}
	preSkip := enc.Lookahead() * 48000 / rate
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := &Writer{f: f, w: bufio.NewWriterSize(f, 64<<10), enc: enc, rate: rate, preSkip: preSkip,
		serial: rand.Uint32(), packet: make([]byte, 4000)}
	head := opusHead(preSkip, rate)
	tags := opusTags("Sonora")
	if err := w.writePage([][]byte{head}, 0, 0x02); err != nil { // beginning of stream
		w.abort()
		return nil, err
	}
	if err := w.writePage([][]byte{tags}, 0, 0); err != nil {
		w.abort()
		return nil, err
	}
	w.granule = uint64(preSkip)
	return w, nil
}

// Write adds samples in [-1, 1] at the writer's rate.
func (w *Writer) Write(samples []float32) error {
	size := w.rate * frameMS / 1000
	for _, s := range samples {
		w.frame = append(w.frame, max(-1, min(1, s)))
		if len(w.frame) == size {
			if err := w.encode(); err != nil {
				return err
			}
		}
	}
	w.written += uint64(len(samples))
	return nil
}

func (w *Writer) encode() error {
	n, err := w.enc.Encode(w.frame, w.packet)
	if err != nil {
		return err
	}
	w.frame = w.frame[:0]
	w.granule += uint64(48000 * frameMS / 1000)
	w.addPacket(w.packet[:n])
	if w.packets >= packetsPerPage || len(w.lacing) > 200 {
		return w.flushPage(0)
	}
	return nil
}

// Samples returns how much audio has been written, in input samples.
func (w *Writer) Samples() uint64 { return w.written }

// Close encodes what's left (padding the last frame with silence),
// marks the end of the stream and closes the file. The last page's
// granule position trims the padding.
func (w *Writer) Close() error {
	if len(w.frame) > 0 {
		size := w.rate * frameMS / 1000
		for len(w.frame) < size {
			w.frame = append(w.frame, 0)
		}
		if err := w.encode(); err != nil {
			w.abort()
			return err
		}
	}
	// Exact end: pre-skip plus the real samples, at 48 kHz.
	w.granule = uint64(w.preSkip) + w.written*48000/uint64(w.rate)
	if w.packets == 0 {
		// An empty end-of-stream page still carries the final position.
		w.lacing = w.lacing[:0]
	}
	if err := w.flushPage(0x04); err != nil { // end of stream
		w.abort()
		return err
	}
	return w.f.Close()
}

func (w *Writer) abort() {
	w.f.Close()
	os.Remove(w.f.Name())
}
