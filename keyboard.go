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
	"sync"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
)

const devicePath = "/dev/input"

// KeyboardDevice represents a physical keyboard, it contains the dev struct, file descriptor and state of any "modifier" keys
type KeyboardDevice struct {
	dev       *C.struct_libevdev
	fd        *os.File
	modifiers *KeyModifiers
}

// Close will gracefully handle closing a keyboard device, freeing memory and file descriptors
func (k *KeyboardDevice) Close() {
	C.libevdev_free(k.dev)
	k.fd.Close()
}

func (k *KeyboardDevice) isKeyboard() bool {
	if C.libevdev_has_event_code(k.dev, C.EV_KEY, C.KEY_CAPSLOCK) == 1 {
		return true
	} else {
		return false
	}
}

// OpenKeyboardDevice will open a specific keyboard device (from the device path passed as a string)
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

// OpenKeyboardDevices will open all currently connected keyboards passing them out through a channel for further processing
func OpenKeyboardDevices() <-chan *KeyboardDevice {
	kbdChan := make(chan *KeyboardDevice)
	var kbdPaths []string
	fileRegexp, _ := regexp.Compile(`event\d+$`)
	log.Debug("Looking for keyboards...")
	err := filepath.WalkDir(devicePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Errorf("could not read %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			if fileRegexp.MatchString(path) {
				kbdPaths = append(kbdPaths, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("Couldn't traverse device path: %s, %v", devicePath, err)
	}
	log.Debug("Keyboard search finished.")
	go func() {
		for _, kbdPath := range kbdPaths {
			kbd := OpenKeyboardDevice(kbdPath)
			if kbd.isKeyboard() {
				log.Debugf("Opening keyboard device %s", kbdPath)
				kbdChan <- kbd
			} else {
				kbd.Close()
			}
		}
		close(kbdChan)
	}()
	return kbdChan
}

// SnoopAllKeyboards will snoop or listen for all key events on all currently connected keyboards.  Keyboards are passed in through a channel, see OpenKeyboardDevices for an example of opening all connected keyboards
func SnoopAllKeyboards(kbds <-chan *KeyboardDevice) <-chan KeyEvent {
	norm := C.enum_libevdev_read_flag(C.LIBEVDEV_READ_FLAG_NORMAL)
	keys := make(chan KeyEvent)
	var wg sync.WaitGroup
	kbdSnoop := func(kbd *KeyboardDevice) {
		defer wg.Done()
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
			e.updateRune(kbd.modifiers)
			keys <- *e
		}
	}
	for kbd := range kbds {
		log.Debugf("Tracking keys on device %s", kbd.fd.Name())
		wg.Add(1)
		go kbdSnoop(kbd)
	}
	go func() {
		defer close(keys)
		wg.Wait()
	}()
	return keys
}

// VirtualKeyboardDevice represents a "virtual" (uinput) keyboard device
type VirtualKeyboardDevice struct {
	uidev *C.struct_libevdev_uinput
	dev   *C.struct_libevdev
}

// NewVirtualKeyboard will create a new virtual keyboard device (with the name passed in)
func NewVirtualKeyboard(name string) *VirtualKeyboardDevice {
	var uidev *C.struct_libevdev_uinput

	dev := C.libevdev_new()
	C.libevdev_set_name(dev, C.CString(name))
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

// SyncEvent sends the EVSYN event to a virtual keyboard, required between other events
func (u *VirtualKeyboardDevice) SyncEvent() error {
	rv := C.libevdev_uinput_write_event(u.uidev, C.EV_SYN, C.SYN_REPORT, 0)
	if rv < 0 {
		return errors.New("failed to issue EV_SYN SYN_REPORT")
	}
	return nil
}

// KeyEvent sends a specific keycode and value to the virtual keyboard
func (u *VirtualKeyboardDevice) KeyEvent(keyCode int, value int) error {
	rv := C.libevdev_uinput_write_event(u.uidev, C.EV_KEY, C.uint(keyCode), C.int(value))
	if rv < 0 {
		name := C.libevdev_event_value_get_name(C.EV_KEY, C.uint(keyCode), 0)
		err := fmt.Errorf("failed to issue key press event (EV_KEY, 0) for %s", C.GoString(name))
		return err
	}
	return u.SyncEvent()
}

// KeyPressEvent sends a specified key "press" to the virtual keyboard
func (u *VirtualKeyboardDevice) KeyPressEvent(keyCode int) error {
	return u.KeyEvent(keyCode, 1)
}

// KeyReleaseEvent sends a specified key "release" to the virtual keyboard
func (u *VirtualKeyboardDevice) KeyReleaseEvent(keyCode int) error {
	return u.KeyEvent(keyCode, 0)
}

// HoldShift sends the equivalent of holding down the shift key to the virtual keyboard
func (u *VirtualKeyboardDevice) HoldShift() error {
	return u.KeyEvent(C.KEY_LEFTSHIFT, 1)
}

// ReleaseShift sends the equivalent of releasing the shift key to the virtual keyboard
func (u *VirtualKeyboardDevice) ReleaseShift() error {
	return u.KeyEvent(C.KEY_LEFTSHIFT, 0)
}

// TypeRune is a high level way to "type" a specific Go rune on the keyboard
func (u *VirtualKeyboardDevice) TypeRune(r rune) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Error in TypeRune operation: %v", r)
		}
	}()
	if !unicode.In(r, unicode.PrintRanges...) {
		err := fmt.Errorf("rune %c (%U) is not a printable character", r, r)
		panic(err)
	}
	keyCode, isUpperCase := CodeAndCase(r)
	if isUpperCase {
		checkErr(u.HoldShift())
	}
	checkErr(u.KeyPressEvent(keyCode))
	checkErr(u.KeyReleaseEvent(keyCode))
	if isUpperCase {
		checkErr(u.ReleaseShift())
	}
}

// TypeSpace is a high level way to "type" a space character (effectively, press/release the spacebar)
func (u *VirtualKeyboardDevice) TypeSpace() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Error in TypeSpace operation: %v", r)
		}
	}()
	checkErr(u.KeyEvent(C.KEY_SPACE, 1))
	checkErr(u.KeyEvent(C.KEY_SPACE, 0))
}

// TypeBackspace allows you to "type" a backspace key and remove a single character
func (u *VirtualKeyboardDevice) TypeBackspace() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Error in TypeBackspace operation: %v", r)
		}
	}()
	checkErr(u.KeyEvent(C.KEY_BACKSPACE, 1))
	checkErr(u.KeyEvent(C.KEY_BACKSPACE, 0))
}

// TypeString is a high level function that makes it easy to "type" out a string to the virtual keyboard
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

// Close will gracefully remove a virtual keyboard, freeing memory and file descriptors
func (u *VirtualKeyboardDevice) Close() {
	C.libevdev_uinput_destroy(u.uidev)
	C.libevdev_free(u.dev)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
