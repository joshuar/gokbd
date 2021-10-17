package gokbd

// #cgo pkg-config: libevdev
// #include <libevdev/libevdev.h>
// #include <libevdev/libevdev-uinput.h>
import "C"
import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
)

const devicePath = "/dev/input"

type KeyboardDevice struct {
	dev       *C.struct_libevdev
	fd        *os.File
	modifiers *KeyModifiers
}

func (k *KeyboardDevice) Close() {
	C.libevdev_free(k.dev)
	k.fd.Close()
}

func (k *KeyboardDevice) IsKeyboard() bool {
	if C.libevdev_has_event_code(k.dev, C.EV_KEY, C.KEY_CAPSLOCK) == 1 {
		return true
	} else {
		return false
	}
}

func OpenKeyboardDevice(devPath string) *KeyboardDevice {
	dev := C.libevdev_new()
	fd, err := os.Open(devPath)
	if err != nil {
		log.Fatalf("Failed to open device: %v", err)
	}
	c_err := C.libevdev_set_fd(dev, C.int(fd.Fd()))
	if c_err > 0 {
		log.Fatalf("Failed to init libevdev: %v", err)
		os.Exit(1)
	}
	return &KeyboardDevice{
		dev:       dev,
		fd:        fd,
		modifiers: NewKeyModifers(),
	}
}

func OpenKeyboardDevices() []*KeyboardDevice {
	var kbds []*KeyboardDevice
	fileRegexp, _ := regexp.Compile(`event\d+$`)
	err := filepath.WalkDir(devicePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Errorf("could not read %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			if fileRegexp.MatchString(path) {
				kbd := OpenKeyboardDevice(path)
				if kbd.IsKeyboard() {
					kbds = append(kbds, kbd)
				} else {
					kbd.Close()
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("Couldn't traverse device path: %s, %v", devicePath, err)
	}
	return kbds
}

func SnoopAllKeyboards(keys chan KeyEvent) error {
	kbds := OpenKeyboardDevices()
	if len(kbds) == 0 {
		return errors.New("no keyboards found to snoop")
	}
	norm := C.enum_libevdev_read_flag(C.LIBEVDEV_READ_FLAG_NORMAL)
	for _, kbd := range kbds {
		log.Debugf("Tracking keys on device %s", kbd.fd.Name())
		go func(kbd *KeyboardDevice) {
			for {
				var ev C.struct_input_event
				C.libevdev_next_event(kbd.dev, C.uint(norm), &ev)
				e := NewKeyEvent(ev)
				if e.Value != 2 {
					switch e.EventName {
					case "KEY_CAPSLOCK":
						kbd.modifiers.ToggleCapsLock()
					case "KEY_LEFTSHIFT", "KEY_RIGHTSHIFT":
						kbd.modifiers.ToggleShift()
					case "KEY_LEFTCTRL", "KEY_RIGHTCTRL":
						kbd.modifiers.ToggleCtrl()
					case "KEY_LEFTALT", "KEY_RIGHTALT":
						kbd.modifiers.ToggleAlt()
					case "KEY_LEFTMETA", "KEY_RIGHTMETA":
						kbd.modifiers.ToggleMeta()
					}
				}
				e.UpdateRune(kbd.modifiers)
				keys <- *e
			}
		}(kbd)
	}
	return nil
}

type VirtualKeyboardDevice struct {
	uidev *C.struct_libevdev_uinput
	dev   *C.struct_libevdev
}

func NewVirtualKeyboard() *VirtualKeyboardDevice {
	var uidev *C.struct_libevdev_uinput

	dev := C.libevdev_new()
	C.libevdev_set_name(dev, C.CString("virtual keyboard"))
	// expose the relevant event types
	C.libevdev_enable_event_type(dev, C.EV_REL)
	C.libevdev_enable_event_type(dev, C.EV_KEY)
	C.libevdev_enable_event_type(dev, C.EV_REP)
	// expose all physical ascii keys on a standard qwerty keyboard
	for k := range runeMap {
		C.libevdev_enable_event_code(dev, C.EV_KEY, C.uint(k), nil)
	}
	// expose some modifier keys (in this case just the left ones, we only need those)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.KEY_LEFTSHIFT, nil)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.KEY_LEFTCTRL, nil)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.KEY_LEFTALT, nil)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.KEY_LEFTMETA, nil)

	rv := C.libevdev_uinput_create_from_device(dev, C.LIBEVDEV_UINPUT_OPEN_MANAGED, &uidev)
	if rv > 0 {
		log.Errorf("Failed to create new uinput device: %v", rv)
		return nil
	}
	log.Debugf("Virtual keyboard created at %s", C.GoString(C.libevdev_uinput_get_devnode(uidev)))
	time.Sleep(1 * time.Second)
	return &VirtualKeyboardDevice{
		uidev: uidev,
		dev:   dev,
	}
}

func (u *VirtualKeyboardDevice) SyncEvent() error {
	rv := C.libevdev_uinput_write_event(u.uidev, C.EV_SYN, C.SYN_REPORT, 0)
	if rv < 0 {
		return errors.New("failed to issue EV_SYN SYN_REPORT")
	}
	return nil
}

func (u *VirtualKeyboardDevice) KeyEvent(keyCode int, value int) error {
	rv := C.libevdev_uinput_write_event(u.uidev, C.EV_KEY, C.uint(keyCode), C.int(value))
	if rv < 0 {
		name := C.libevdev_event_value_get_name(C.EV_KEY, C.uint(keyCode), 0)
		err := fmt.Errorf("failed to issue key press event (EV_KEY, 0) for %s", C.GoString(name))
		return err
	}
	return u.SyncEvent()
}

func (u *VirtualKeyboardDevice) KeyPressEvent(keyCode int) error {
	return u.KeyEvent(keyCode, 1)
}

func (u *VirtualKeyboardDevice) KeyReleaseEvent(keyCode int) error {
	return u.KeyEvent(keyCode, 0)
}

func (u *VirtualKeyboardDevice) HoldShift() error {
	return u.KeyEvent(C.KEY_LEFTSHIFT, 1)
}

func (u *VirtualKeyboardDevice) ReleaseShift() error {
	return u.KeyEvent(C.KEY_LEFTSHIFT, 0)
}

func (u *VirtualKeyboardDevice) TypeRune(r rune) {
	if !unicode.In(r, unicode.PrintRanges...) {
		err := fmt.Errorf("rune %c (%U) is not a printable character", r, r)
		log.Error(err)
	}
	keyCode, isUpperCase := CodeAndCase(r)
	if isUpperCase {
		u.HoldShift()
		u.KeyPressEvent(keyCode)
		u.KeyReleaseEvent(keyCode)
		u.ReleaseShift()
	} else {
		u.KeyPressEvent(keyCode)
		u.KeyReleaseEvent(keyCode)
	}
}

func (u *VirtualKeyboardDevice) TypeSpace() {
	u.KeyEvent(C.KEY_SPACE, 1)
	u.KeyEvent(C.KEY_SPACE, 0)
}

func (u *VirtualKeyboardDevice) TypeBackspace() {
	u.KeyEvent(C.KEY_BACKSPACE, 1)
	u.KeyEvent(C.KEY_BACKSPACE, 0)
}

func (u *VirtualKeyboardDevice) TypeString(str string) {
	s := strings.NewReader(str)
	for {
		r, _, err := s.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("Error reading rune in string: %v", err)
		}
		switch r {
		case ' ':
			u.TypeSpace()
		default:
			u.TypeRune(r)
		}
	}
}

func (u *VirtualKeyboardDevice) Close() {
	C.libevdev_uinput_destroy(u.uidev)
	C.libevdev_free(u.dev)
}
