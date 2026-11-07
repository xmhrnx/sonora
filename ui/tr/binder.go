package tr

import (
	"sync"

	"github.com/xmhrnx/giode/i18n"
)

var (
	mu    sync.Mutex
	binds []func()
	once  sync.Once
)

// Bind sets an element's text to T(key, args...) now and again after
// every language change. Use it for elements built once; text drawn
// every frame calls T instead.
//
//	tr.Bind(func(s string) { b.Label(s) }, "Record")
func Bind(set func(string), key string, args ...any) {
	once.Do(func() { i18n.OnChange(rebind) })
	f := func() { set(T(key, args...)) }
	f()
	mu.Lock()
	binds = append(binds, f)
	mu.Unlock()
}

func rebind() {
	mu.Lock()
	fs := append([]func(){}, binds...)
	mu.Unlock()
	for _, f := range fs {
		f()
	}
}

// Msg marks text that's translated later, where it's shown (an error
// message, a status): it returns key unchanged, and the catalog check
// finds it.
func Msg(key string) string { return key }
