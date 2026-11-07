package library

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mehran9675/sonora/internal/audio/opusfile"
	"github.com/mehran9675/sonora/ui/core"
)

// recording is an audio file in the recordings folder.
type recording struct {
	Path     string
	Name     string
	Modified time.Time
	Size     int64
	Duration time.Duration // 0 when unknown
}

// recordingList scans the recordings folder, at most once a second.
type recordingList struct {
	c       *core.Ctx
	scanned time.Time
	dir     string
	items   []recording
	lengths map[string]time.Duration // by path and modification time
}

// audioExts are the files listed as recordings.
var audioExts = map[string]bool{".opus": true, ".ogg": true, ".wav": true, ".mp3": true, ".flac": true, ".m4a": true}

func (l *recordingList) list() []recording {
	dir := l.c.RecordingsDir()
	if dir == l.dir && time.Since(l.scanned) < time.Second {
		return l.items
	}
	if l.lengths == nil {
		l.lengths = map[string]time.Duration{}
	}
	l.dir, l.scanned, l.items = dir, time.Now(), nil
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		info, err := e.Info()
		if e.IsDir() || !audioExts[ext] || err != nil {
			continue
		}
		r := recording{Path: filepath.Join(dir, e.Name()), Name: strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())),
			Modified: info.ModTime(), Size: info.Size()}
		key := r.Path + "|" + info.ModTime().String()
		if d, ok := l.lengths[key]; ok {
			r.Duration = d
		} else if ext == ".opus" || ext == ".ogg" {
			r.Duration, _ = opusfile.Duration(r.Path)
			l.lengths[key] = r.Duration
		}
		l.items = append(l.items, r)
	}
	sort.Slice(l.items, func(i, j int) bool { return l.items[i].Modified.After(l.items[j].Modified) })
	return l.items
}

// refresh rescans on the next frame.
func (l *recordingList) refresh() { l.scanned = time.Time{} }
