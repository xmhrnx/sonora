package settings

import (
	"runtime"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/convert"
	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/pagehead"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/tr"
)

func (p *Page) content() giode.Element {
	toolchain := giode.Element(chip.Accent(p.uvInfo, "theme.success"))
	if p.uvInfo == "" {
		toolchain = giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter},
			chip.Accent(tr.T("Downloads on first install"), "theme.warning"), p.uv)
	}
	ffmpeg := chip.Accent(tr.T("Found"), "theme.success")
	if !convert.Available() {
		ffmpeg = chip.Accent(tr.T("Not found — needed for M4A, Opus and video files"), "theme.muted")
	}
	sections := []giode.Element{
		pagehead.Head(giode.Compact(), tr.T("Settings"), tr.T("Make Sonora yours.")),
		section(tr.T("Appearance"), tr.T("How Sonora looks."),
			p.setting(tr.T("Theme"), tr.T("System follows your device's light or dark setting."), p.theme.Element()),
			p.setting(tr.T("Language"), tr.T("System uses your device's language when Sonora has it."), p.language.trigger.Element()),
		),
		p.recordingSection(),
	}
	// Phones run only on-device models, so the compute device, folders
	// and the Python toolchain don't apply there.
	if !core.Mobile {
		sections = append(sections,
			section(tr.T("Performance"), tr.T("Where models run. GPUs are much faster for large models."),
				p.setting(tr.T("Compute device"), tr.T("Automatic uses CUDA or Apple silicon when available."), p.device.Element()),
				p.setting(tr.T("Hugging Face token"), tr.T("Only for gated or private models."), giode.Box(giode.Styles{Width: -1}, p.token)),
			),
			section(tr.T("Storage"), tr.T("Engines, models and your clips live here."),
				p.setting(tr.T("Data folder"), paths.Data(), p.data),
				p.setting(tr.T("Library"), p.c.Library.Dir(), p.library),
			),
			section(tr.T("Toolchain"), tr.T("Sonora manages Python for you; nothing to install by hand."),
				p.setting(tr.T("Python manager (uv)"), tr.T("Creates an isolated environment per engine."), toolchain),
				p.setting(tr.T("FFmpeg"), tr.T("Reads audio formats Sonora can't decode itself."), ffmpeg),
			),
		)
		if row := p.fileManager.element(p); row != nil {
			sections = append(sections, section(tr.T("Desktop"), tr.T("How Sonora fits into your computer."), row))
		}
	}
	sections = append(sections,
		section(tr.T("About"), "Sonora "+Version+" · "+runtime.GOOS+"/"+runtime.GOARCH,
			giode.Text(tr.T("Built with giode, a Go UI library on Gio. Engines are open-source projects by their authors; "+
				"check each model's license before commercial use."), giode.Styles{Color: "theme.secondary", FontSize: 13, LineHeight: 20}),
		),
	)
	return giode.Stack(giode.Styles{Width: -1, Gap: 20, Padding: giode.UniformInset(core.Gutter()), MaxWidth: 900}, sections...)
}

// Element lays out the page.
func (p *Page) Element() giode.Element {
	return giode.Box(giode.Styles{Width: -1, Height: -1}, p.scroll)
}
