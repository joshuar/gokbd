package gokbd

// #cgo pkg-config: libevdev
// #include <libevdev/libevdev.h>
// #include <libevdev/libevdev-uinput.h>
import "C"

type KeyEvent struct {
	eventRaw  C.struct_input_event
	Value     int
	TypeName  string
	EventName string
	AsRune    rune
}

func NewKeyEvent(ev C.struct_input_event) *KeyEvent {
	return &KeyEvent{
		eventRaw:  ev,
		Value:     int(ev.value),
		TypeName:  C.GoString(C.libevdev_event_type_get_name(C.uint(ev._type))),
		EventName: C.GoString(C.libevdev_event_code_get_name(C.uint(ev._type), C.uint(ev.code))),
		AsRune:    runeMap[int(ev.code)].lc,
	}
}

func (kev *KeyEvent) UpdateRune(modifiers *KeyModifiers) {
	switch {
	case modifiers.CapsLock:
	case modifiers.Shift:
		kev.AsRune = runeMap[int(kev.eventRaw.code)].uc
	}
}

func (kev *KeyEvent) IsKeyPress() bool {
	if kev.Value == 1 && kev.TypeName == "EV_KEY" {
		return true
	}
	return false
}

func (kev *KeyEvent) IsKeyRelease() bool {
	if kev.Value == 0 && kev.TypeName == "EV_KEY" {
		return true
	}
	return false
}

func (kev *KeyEvent) IsBackspace() bool {
	switch kev.EventName {
	case "KEY_BACKSPACE":
		return true
	default:
		return false
	}
}

func (kev *KeyEvent) IsModifier() bool {
	switch kev.EventName {
	case "KEY_LEFTCTRL", "KEY_RIGHTCTRL", "KEY_LEFTALT", "KEY_RIGHTALT", "KEY_LEFTMETA", "KEY_RIGHTMETA", "KEY_LEFTSHIFT", "KEY_RIGHTSHIFT":
		return true
	default:
		return false
	}
}
