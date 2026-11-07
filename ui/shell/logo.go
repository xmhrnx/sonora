package shell

import (
	"bytes"
	"image"
	"image/png"
	"sync"

	"github.com/xmhrnx/giode"
	"golang.org/x/image/draw"

	"github.com/mehran9675/sonora/assets"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// logoImage is the app icon, decoded once and shrunk for the sidebar so
// it stays crisp.
var logoImage = sync.OnceValue(func() image.Image {
	src, err := png.Decode(bytes.NewReader(assets.Icon))
	if err != nil {
		return image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	dst := image.NewNRGBA(image.Rect(0, 0, 128, 128))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return dst
})

// logo is the app icon and name at the top of the sidebar.
func logo() giode.Element {
	return giode.Row(giode.Styles{Gap: 12, Align: giode.AlignCenter, Padding: giode.SymmetricInset(6, 0), Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{Width: 38, Height: 38, BorderRadius: 10,
			BoxShadow: giode.BoxShadow{Y: 6, Blur: 18, Color: "#8b5cf644"}},
			giode.Image(logoImage(), giode.Styles{Width: 38, Height: 38, Fit: giode.FitContain})),
		giode.Stack(giode.Styles{Gap: 1},
			giode.Text("Sonora", giode.Styles{Color: "theme.text", FontSize: 18, FontWeight: giode.FontWeightExtraBold}),
			giode.Text(tr.T("Speech studio"), theme.Caption),
		),
	)
}
