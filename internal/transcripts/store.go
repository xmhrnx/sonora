package transcripts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Open loads the store in dir.
func Open(dir string) *Store {
	s := &Store{dir: dir}
	os.MkdirAll(dir, 0o755)
	if data, err := os.ReadFile(s.index()); err == nil {
		json.Unmarshal(data, &s.entries)
	}
	return s
}

// NewFile returns a fresh id and the path its Opus audio should be
// written to.
func (s *Store) NewFile() (id, path string) {
	id = fmt.Sprintf("%d", time.Now().UnixNano())
	return id, filepath.Join(s.dir, id+".opus")
}

// AudioPath returns the absolute path of an entry's audio.
func (s *Store) AudioPath(e Entry) string { return filepath.Join(s.dir, e.Audio) }

// Add records an entry, newest first.
func (s *Store) Add(e Entry) {
	s.mu.Lock()
	s.entries = append([]Entry{e}, s.entries...)
	s.mu.Unlock()
	s.save()
}

// SetText replaces an entry's text (after the user edits it).
func (s *Store) SetText(id, text string) {
	s.mu.Lock()
	for i := range s.entries {
		if s.entries[i].ID == id {
			s.entries[i].Text = text
		}
	}
	s.mu.Unlock()
	s.save()
}

// Get returns the entry with id.
func (s *Store) Get(id string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.entries {
		if e.ID == id {
			return e, true
		}
	}
	return Entry{}, false
}

// Remove deletes an entry and its audio.
func (s *Store) Remove(id string) {
	s.mu.Lock()
	for i, e := range s.entries {
		if e.ID == id {
			os.Remove(filepath.Join(s.dir, e.Audio))
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
	s.save()
}

// Entries returns a copy of all entries, newest first.
func (s *Store) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := append([]Entry(nil), s.entries...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Created.After(out[j].Created) })
	return out
}

func (s *Store) save() {
	s.mu.Lock()
	data, _ := json.MarshalIndent(s.entries, "", "  ")
	s.mu.Unlock()
	tmp := s.index() + ".tmp"
	if os.WriteFile(tmp, data, 0o644) == nil {
		os.Rename(tmp, s.index())
	}
}
