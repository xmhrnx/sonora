package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

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

// Open reads the profiles in dir; a missing or unreadable index starts
// empty.
func Open(dir string) *Store {
	s := &Store{dir: dir}
	if data, err := os.ReadFile(s.index()); err == nil {
		json.Unmarshal(data, &s.list)
	}
	return s
}
