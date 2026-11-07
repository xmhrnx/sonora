package settings

import (
	"context"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/internal/pyenv"
	st "github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/segmented"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Page is the Settings screen. Create it once.
type Page struct {
	c       *core.Ctx
	scroll  *scroll.Scroll
	theme   *segmented.Control
	device  *segmented.Control
	token   *input.Input
	data    *giode.ButtonElement
	library *giode.ButtonElement
	uv      *giode.ButtonElement
	uvInfo  string

	quality, limit       *picker.Picker
	qualityT, limitT     *picker.Trigger
	folder               *input.Input
	openFolder, defaults *giode.ButtonElement
	mic                  microphone
	language             language
	fileManager          fileManager
}

// Version is set at build time with -ldflags "-X …/settings.Version=…".
var Version = "dev"

// New builds the page; applyTheme switches the palette.
func New(c *core.Ctx, applyTheme func(mode string)) *Page {
	p := &Page{c: c}
	if !core.Mobile {
		p.uvInfo = pyenv.UVVersion()
	}
	p.theme = segmented.New([]segmented.Option{{Value: "system", Label: tr.Msg("System")}, {Value: "dark", Label: tr.Msg("Dark")}, {Value: "light", Label: tr.Msg("Light")}},
		func() string { return c.Settings.Get().Theme }, applyTheme)
	p.device = segmented.New([]segmented.Option{{Value: "auto", Label: tr.Msg("Automatic")}, {Value: "gpu", Label: tr.Msg("GPU")}, {Value: "cpu", Label: tr.Msg("CPU only")}},
		func() string { return c.Settings.Get().Device },
		func(v string) {
			c.Settings.Update(func(s *st.Settings) { s.Device = v })
			c.Manager.CloseIdle() // engines restart on the new device on next use
			c.Info(map[string]string{
				"auto": tr.T("Engines will restart and pick the device themselves"),
				"gpu":  tr.T("Engines will restart on the GPU"),
				"cpu":  tr.T("Engines will restart on the CPU"),
			}[v])
		})
	p.token = giode.Input.Text("hf_…", theme.Field()).Password(true).
		OnChange(func(v string) { c.Settings.Update(func(s *st.Settings) { s.HFToken = v }) })
	tr.Bind(func(s string) { p.token.Label(s) }, "Hugging Face token")
	p.token.SetText(c.Settings.Get().HFToken)
	outline := giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border", FontSize: 13,
		Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}
	p.data = giode.Button("", outline).Variant(giode.VariantOutlined).
		Icon(icons.Icon(icons.Folder, 15, "")).OnClick(func() { c.Error(paths.Open(paths.Data())) })
	tr.Bind(func(s string) { p.data.Label(s) }, "Open data folder")
	p.library = giode.Button("", outline).Variant(giode.VariantOutlined).
		Icon(icons.Icon(icons.Folder, 15, "")).OnClick(func() { c.Error(paths.Open(c.Library.Dir())) })
	tr.Bind(func(s string) { p.library.Label(s) }, "Open library folder")
	p.uv = giode.Button("", outline).Variant(giode.VariantOutlined).Icon(icons.Icon(icons.Download, 15, ""))
	tr.Bind(func(s string) { p.uv.Label(s) }, "Set up now")
	p.uv.OnClick(func() {
		p.uv.Loading(true)
		c.Background(func() error {
			_, err := pyenv.EnsureUV(context.Background(), func(string) {})
			return err
		}, func(err error) {
			p.uv.Loading(false)
			if err != nil {
				c.Error(err)
				return
			}
			p.uvInfo = pyenv.UVVersion()
			c.Success(tr.T("Python toolchain ready"))
		})
	})
	p.recordingControls()
	p.mic.init(c)
	p.language.init(c)
	p.fileManager.init(c)
	p.scroll = giode.Scroll(p.content, theme.Scrolling(6))
	return p
}
