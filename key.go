// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

// #cgo pkg-config: libevdev
// #include <libevdev/libevdev.h>
// #include <libevdev/libevdev-uinput.h>
import "C"

type key struct {
	keyType, keyCode, value int
}

func keyPress(c int) *key {
	return &key{
		keyType: C.EV_KEY,
		keyCode: c,
		value:   1,
	}
}

func keyRelease(c int) *key {
	return &key{
		keyType: C.EV_KEY,
		keyCode: c,
		value:   0,
	}
}

func keySync() *key {
	return &key{
		keyType: C.EV_SYN,
		keyCode: C.SYN_REPORT,
		value:   0,
	}
}

func keySequence(keys ...*key) <-chan *key {
	out := make(chan *key)
	go func() {
		for _, n := range keys {
			out <- n
		}
		close(out)
	}()
	return out
}
