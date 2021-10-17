package gokbd

// KeyModifiers represents the state of any "modifier" keys on the keyboard
type KeyModifiers struct {
	CapsLock bool
	Alt      bool
	Ctrl     bool
	Shift    bool
	Meta     bool
}

// ToggleAlt keeps track of whether an Alt key has been pressed
func (km *KeyModifiers) ToggleAlt() {
	km.Alt = !km.Alt
}

// ToggleShift keeps track of whether an Shift key has been pressed
func (km *KeyModifiers) ToggleShift() {
	km.Shift = !km.Shift
}

// ToggleCtrl keeps track of whether an Ctrl key has been pressed
func (km *KeyModifiers) ToggleCtrl() {
	km.Ctrl = !km.Ctrl
}

// ToggleMeta keeps track of whether an Meta key has been pressed
func (km *KeyModifiers) ToggleMeta() {
	km.Meta = !km.Meta
}

// ToggleCapsLock keeps track of whether the Caps Lock key has been pressed
func (km *KeyModifiers) ToggleCapsLock() {
	km.CapsLock = !km.CapsLock
}

// NewKeyModifiers sets up a struct for tracking whether any of the modifier
// keys have been pressed
func NewKeyModifers() *KeyModifiers {
	return &KeyModifiers{
		CapsLock: false,
		Alt:      false,
		Ctrl:     false,
		Shift:    false,
		Meta:     false,
	}
}
