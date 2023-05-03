// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

// #cgo pkg-config: libevdev
// #include <libevdev/libevdev.h>
// #include <libevdev/libevdev-uinput.h>
import "C"

// KeyEvent represents an event received from the keyboard
// eventRaw is the libevdev input event, see https://www.kernel.org/doc/html/v4.17/input/input.html#event-interface
// Value is the event value, for example 1 for key press, 0 for key release
// TypeName is the event type as a string, for example EV_KEY or EV_SYN
// EventName is the event name as a string, for example KEY_A
// AsRune is the key as a Go rune, for example 'a'
type KeyEvent struct {
	eventRaw  C.struct_input_event
	Value     int
	TypeName  string
	EventName string
	AsRune    rune
}

// NewKeyEvent will create a new key event for whatever just happened on the keyboard
func NewKeyEvent(ev C.struct_input_event) *KeyEvent {
	return &KeyEvent{
		eventRaw:  ev,
		Value:     int(ev.value),
		TypeName:  C.GoString(C.libevdev_event_type_get_name(C.uint(ev._type))),
		EventName: C.GoString(C.libevdev_event_code_get_name(C.uint(ev._type), C.uint(ev.code))),
		AsRune:    runeMap[int(ev.code)].lc,
	}
}

func (kev *KeyEvent) updateRune(modifiers *KeyModifiers) {
	switch {
	case modifiers.CapsLock:
		fallthrough
	case modifiers.Shift:
		kev.AsRune = runeMap[int(kev.eventRaw.code)].uc
	}
}

// IsKeyPress will return true when the event represents a key being pressed
func (kev *KeyEvent) IsKeyPress() bool {
	if kev.Value == 1 && kev.TypeName == "EV_KEY" {
		return true
	}
	return false
}

// IsKeyRelease will return true when the event represents a key being released
func (kev *KeyEvent) IsKeyRelease() bool {
	if kev.Value == 0 && kev.TypeName == "EV_KEY" {
		return true
	}
	return false
}

// IsBackspace will return true when the event represents a key involved is the backspace key
func (kev *KeyEvent) IsBackspace() bool {
	switch kev.EventName {
	case "KEY_BACKSPACE":
		return true
	default:
		return false
	}
}

// IsModifier will return true when the event represents any of the "modifier" keys: Ctrl, Alt, Meta or Shift
func (kev *KeyEvent) IsModifier() bool {
	switch kev.EventName {
	case "KEY_LEFTCTRL", "KEY_RIGHTCTRL", "KEY_LEFTALT", "KEY_RIGHTALT", "KEY_LEFTMETA", "KEY_RIGHTMETA", "KEY_LEFTSHIFT", "KEY_RIGHTSHIFT":
		return true
	default:
		return false
	}
}
