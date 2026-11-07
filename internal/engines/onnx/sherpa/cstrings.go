//go:build sherpa || android

package sherpa

/*
#include <stdlib.h>
#include "c-api.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// cstrings allocates C strings and frees them together.
type cstrings []*C.char

func (s *cstrings) str(v string) *C.char {
	if v == "" {
		return nil
	}
	c := C.CString(v)
	*s = append(*s, c)
	return c
}

func (s cstrings) free() {
	for _, c := range s {
		C.free(unsafe.Pointer(c))
	}
}

func threads(n int) C.int32_t {
	if n <= 0 {
		n = max(1, min(4, runtime.NumCPU()))
	}
	return C.int32_t(n)
}
