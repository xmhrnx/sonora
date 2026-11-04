// Package profiles keeps voice profiles: voices to clone, each a short
// recording of one person and what they say in it, so the engine gets
// both without the user writing anything.
package profiles

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/audio/opusfile"
)

// Profile is a voice to clone.
type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Audio is the recording inside the folder, Opus.
	Audio string `json:"audio"`
	// Transcript is what is said in the recording: the passage read
	// aloud, or what a transcription engine heard in an imported file.
	// Empty until known.
	Transcript string        `json:"transcript,omitempty"`
	Duration   time.Duration `json:"duration"`
	// Imported marks profiles made from a file rather than recorded.
	Imported bool      `json:"imported,omitempty"`
	Created  time.Time `json:"created"`
}

// MaxLength is the longest recording a profile keeps; engines clone
// from 5-30 seconds and longer imports are cut.
const MaxLength = 30 * time.Second

// bitrate keeps voices clear enough to clone.
const bitrate = 64000

// rate is the rate profiles are stored at.
const rate = 48000

// Store is the folder of profiles plus index.json. Methods are safe
// from any goroutine.
type Store struct {
	mu   sync.Mutex
	dir  string
	list []Profile
}

// Open reads the profiles in dir; a missing or unreadable index starts
// empty.
func Open(dir string) *Store {
	s := &Store{dir: dir}
	if data, err := os.ReadFile(s.index()); err == nil {
		json.Unmarshal(data, &s.list)
	}
	return s
}

func (s *Store) index() string { return filepath.Join(s.dir, "index.json") }

// List returns the profiles, newest first.
func (s *Store) List() []Profile {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := append([]Profile(nil), s.list...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Created.After(out[j].Created) })
	return out
}

// Get returns the profile with id.
func (s *Store) Get(id string) (Profile, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.list {
		if p.ID == id {
			return p, true
		}
	}
	return Profile{}, false
}

// Path is where p's recording is.
func (s *Store) Path(p Profile) string { return filepath.Join(s.dir, p.Audio) }

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
