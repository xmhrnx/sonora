package profiles

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mehran9675/sonora/internal/audio"
	"github.com/mehran9675/sonora/internal/audio/opusfile"
	"github.com/mehran9675/sonora/internal/audio/wav"
)

// tone writes secs seconds of a 220 Hz tone as a WAV.
func tone(t *testing.T, secs float64) string {
	const rate = 16000
	clip := &audio.Clip{SampleRate: rate, Samples: make([]float32, int(secs*rate))}
	for i := range clip.Samples {
		clip.Samples[i] = float32(0.3 * math.Sin(2*math.Pi*220*float64(i)/rate))
	}
	path := filepath.Join(t.TempDir(), "tone.wav")
	if err := wav.Write(path, clip); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestStore(t *testing.T) {
	dir := t.TempDir()
	s := Open(dir)
	long, err := s.Add("Long", "", tone(t, 42), true)
	if err != nil {
		t.Fatal(err)
	}
	if long.Duration != MaxLength {
		t.Errorf("an import of 42 s keeps %v, want %v", long.Duration, MaxLength)
	}
	if d, err := opusfile.Duration(s.Path(long)); err != nil || d.Round(100*time.Millisecond) != MaxLength {
		t.Errorf("the stored file lasts %v (%v)", d, err)
	}
	time.Sleep(time.Millisecond) // distinct creation times
	short, err := s.Add("Mine", Passages[0], tone(t, 6), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add("Blip", "", tone(t, 0.3), false); err == nil {
		t.Error("a 0.3 s recording made a profile")
	}

	// A fresh store reads the index back, newest first.
	got := Open(dir).List()
	if len(got) != 2 || got[0].ID != short.ID || got[1].ID != long.ID || got[0].Transcript != Passages[0] {
		t.Fatalf("listed %+v", got)
	}

	if err := s.Update(long.ID, func(p *Profile) { p.Name, p.Transcript = "Renamed", "hello" }); err != nil {
		t.Fatal(err)
	}
	if p, _ := Open(dir).Get(long.ID); p.Name != "Renamed" || p.Transcript != "hello" {
		t.Errorf("after Update: %+v", p)
	}
	if err := s.Delete(long.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(s.Path(long)); !os.IsNotExist(err) {
		t.Error("Delete left the recording")
	}
	if len(Open(dir).List()) != 1 {
		t.Error("Delete left the profile in the index")
	}
}

// Every passage reads in about 15-25 seconds: long enough to clone
// from, short enough for the 30 s limit.
func TestPassages(t *testing.T) {
	for i, p := range Passages {
		words := len(strings.Fields(p))
		if words < 30 || words > 65 {
			t.Errorf("passage %d has %d words", i, words)
		}
	}
}
