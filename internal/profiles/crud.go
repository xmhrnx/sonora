package profiles

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/audio/opusfile"
)

// Add makes a profile named name from the recording at src, in any
// format decode reads: kept as Opus, cut to MaxLength. src itself is
// left alone. It decodes and encodes, so call it off the UI goroutine.
func (s *Store) Add(name, transcript, src string, imported bool) (Profile, error) {
	clip, err := decode.File(src, rate)
	if err != nil {
		return Profile{}, err
	}
	if max := int(MaxLength.Seconds()) * rate; len(clip.Samples) > max {
		clip.Samples = clip.Samples[:max]
	}
	if clip.Duration() < time.Second {
		return Profile{}, errors.New("the recording is too short")
	}
	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	p := Profile{ID: id, Name: name, Audio: id + ".opus", Transcript: transcript,
		Duration: clip.Duration(), Imported: imported, Created: time.Now()}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Profile{}, err
	}
	w, err := opusfile.Create(s.Path(p), rate, bitrate)
	if err != nil {
		return Profile{}, err
	}
	if err := w.Write(clip.Samples); err != nil {
		w.Close()
		os.Remove(s.Path(p))
		return Profile{}, err
	}
	if err := w.Close(); err != nil {
		os.Remove(s.Path(p))
		return Profile{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list = append(s.list, p)
	return p, s.save()
}

// Update changes the profile with id (its name or transcript).
func (s *Store) Update(id string, fn func(p *Profile)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.list {
		if s.list[i].ID == id {
			fn(&s.list[i])
			return s.save()
		}
	}
	return errors.New("the profile no longer exists")
}

// Delete removes the profile with id and its recording.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.list {
		if p.ID == id {
			os.Remove(s.Path(p))
			s.list = append(s.list[:i], s.list[i+1:]...)
			return s.save()
		}
	}
	return nil
}

// save writes the index; call with mu held.
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.list, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.index() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.index())
}
