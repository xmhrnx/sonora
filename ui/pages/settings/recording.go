package settings

import (
	"os"
	"strconv"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/paths"
	st "github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// recordingControls builds the Recording section's controls.
func (p *Page) recordingControls() {
	c := p.c
	quality := func() string {
		q := c.Settings.Get().RecordQuality
		if q == "" {
			return core.RecordQualities[1].Value
		}
		return q
	}
	p.quality = picker.New(tr.Msg("Recording quality"), func() []picker.Item {
		var items []picker.Item
		for _, q := range core.RecordQualities {
			items = append(items, picker.Item{ID: q.Value, Label: tr.T(q.Name), Detail: q.Detail()})
		}
		return items
	}, quality, func(it picker.Item) { c.Settings.Update(func(s *st.Settings) { s.RecordQuality = it.ID }) })
	p.qualityT = picker.NewTrigger(tr.Msg("Recording quality"), func() (string, string) {
		for _, q := range core.RecordQualities {
			if q.Value == quality() {
				return tr.T(q.Name), q.Detail()
			}
		}
		return quality(), ""
	}, nil, p.quality.Open)
	minutes := func() string { return strconv.Itoa(int(c.RecordMax().Minutes())) }
	p.limit = picker.New(tr.Msg("Longest recording"), func() []picker.Item {
		var items []picker.Item
		for _, m := range core.RecordLimits {
			items = append(items, picker.Item{ID: strconv.Itoa(m), Label: limitLabel(m)})
		}
		return items
	}, minutes, func(it picker.Item) {
		m, _ := strconv.Atoi(it.ID)
		c.Settings.Update(func(s *st.Settings) { s.RecordMinutes = m })
	})
	p.limitT = picker.NewTrigger(tr.Msg("Longest recording"), func() (string, string) {
		m, _ := strconv.Atoi(minutes())
		return limitLabel(m), tr.T("Then the recording stops by itself")
	}, nil, p.limit.Open)
	p.folder = giode.Input.Text(core.DefaultRecordingsDir(), theme.Field()).
		OnChange(func(v string) { c.Settings.Update(func(s *st.Settings) { s.RecordDir = strings.TrimSpace(v) }) })
	tr.Bind(func(s string) { p.folder.Label(s) }, "Recordings folder")
	p.folder.SetText(c.Settings.Get().RecordDir)
	outline := giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border", FontSize: 13,
		Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}
	p.openFolder = giode.Button("", outline).Variant(giode.VariantOutlined).Icon(icons.Icon(icons.Folder, 14, "")).
		OnClick(func() {
			dir := c.RecordingsDir()
			if err := os.MkdirAll(dir, 0o755); err != nil {
				c.Error(err)
				return
			}
			c.Error(paths.Open(dir))
		})
	tr.Bind(func(s string) { p.openFolder.Label(s) }, "Open")
	p.defaults = giode.Button("", outline).Variant(giode.VariantOutlined).OnClick(func() {
		p.folder.SetText("")
		c.Settings.Update(func(s *st.Settings) { s.RecordDir = "" })
	})
	tr.Bind(func(s string) { p.defaults.Label(s) }, "Default")
}

// limitLabel names a recording length in minutes.
func limitLabel(m int) string {
	switch {
	case m == 60:
		return tr.T("1 hour")
	case m%60 == 0:
		return tr.N("{n} hours", m/60)
	}
	return tr.N("{n} minutes", m)
}

// recordingSection is the Recording settings card.
func (p *Page) recordingSection() giode.Element {
	rows := []giode.Element{
		p.setting(tr.T("Microphone"), tr.T("Where recordings come from. An unplugged microphone falls back to the default."),
			p.mic.trigger.Element()),
		p.setting(tr.T("Quality"), tr.T("Recordings are Opus files, a fraction of a WAV's size. Higher quality suits voices you'll clone."),
			p.qualityT.Element()),
		p.setting(tr.T("Longest recording"), tr.T("A safety stop, so a forgotten recording doesn't fill the storage."), p.limitT.Element()),
	}
	// Phones keep recordings in Sonora's own storage; export them from
	// the Library.
	if !core.Mobile {
		rows = append(rows, p.setting(tr.T("Folder"), tr.T("Where recordings are saved. Leave empty for Sonora's data folder."),
			giode.Stack(giode.Styles{Width: -1, Gap: 8},
				giode.Box(giode.Styles{Width: -1}, p.folder),
				giode.Row(giode.Styles{Gap: 8}, p.openFolder, p.defaults),
			)))
	}
	return section(tr.T("Recording"), tr.T("Voices you record to clone or to transcribe."), rows...)
}
