package tr

import (
	"strings"

	"github.com/xmhrnx/giode/i18n"
)

// Language is a language Sonora is translated into.
type Language struct {
	Code string // "de", "pt", "nb"
	Name string // its own name for itself: "Deutsch"
}

// Languages are the translations, English first, then by name.
var Languages = []Language{
	{"en", "English"},
	{"bg", "Български"},
	{"ca", "Català"},
	{"cs", "Čeština"},
	{"da", "Dansk"},
	{"de", "Deutsch"},
	{"et", "Eesti"},
	{"el", "Ελληνικά"},
	{"es", "Español"},
	{"fr", "Français"},
	{"hr", "Hrvatski"},
	{"it", "Italiano"},
	{"lv", "Latviešu"},
	{"lt", "Lietuvių"},
	{"hu", "Magyar"},
	{"nl", "Nederlands"},
	{"nb", "Norsk bokmål"},
	{"pl", "Polski"},
	{"pt", "Português"},
	{"ro", "Română"},
	{"ru", "Русский"},
	{"sk", "Slovenčina"},
	{"sl", "Slovenščina"},
	{"fi", "Suomi"},
	{"sv", "Svenska"},
	{"tr", "Türkçe"},
	{"uk", "Українська"},
}

// Use switches to the language code, or the system's when code is "".
// A system language Sonora isn't translated into reads in English.
func Use(code string) {
	if code == "" {
		code = System()
	}
	i18n.SetLocale(code)
}

// System returns the system's language as one of Languages' codes, or
// "en".
func System() string {
	loc := i18n.SystemLocale()
	lang, _, _ := strings.Cut(loc, "-")
	switch lang {
	case "no", "nn": // Norwegian: Sonora has Bokmål
		lang = "nb"
	}
	for _, l := range Languages {
		if l.Code == lang {
			return lang
		}
	}
	return "en"
}

// Current returns the language in use.
func Current() string {
	lang, _, _ := strings.Cut(i18n.Locale(), "-")
	return lang
}

// Name returns a language's own name.
func Name(code string) string {
	for _, l := range Languages {
		if l.Code == code {
			return l.Name
		}
	}
	return code
}

// T returns key's translation with its {placeholders} filled from
// name, value pairs.
func T(key string, args ...any) string { return i18n.T(key, args...) }

// N returns the plural form of key for n; {n} is n.
func N(key string, n int, args ...any) string { return i18n.N(key, n, args...) }
