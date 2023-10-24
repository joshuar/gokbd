// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

// #cgo pkg-config: libevdev
// #include <libevdev/libevdev.h>
// #include <libevdev/libevdev-uinput.h>
import "C"
import (
	"context"
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

	"github.com/rs/zerolog/log"
	"kernel.org/pub/linux/libs/security/libcap/cap"
)

const devicePath = "/dev/input"

func grabDevice(dev *C.struct_libevdev) (func() error, error) {
	rv := C.libevdev_grab(dev, C.LIBEVDEV_GRAB)
	if rv < 0 {
		return nil, errors.New("failed to grab device")
	}
	ungrab := func() error {
		rv := C.libevdev_grab(dev, C.LIBEVDEV_UNGRAB)
		if rv < 0 {
			return errors.New("failed to ungrab device")
		}
		return nil
	}
	return ungrab, nil
}

// KeyboardDevice represents a physical keyboard, it contains the dev struct,
// file descriptor and state of any "modifier" keys
type KeyboardDevice struct {
	dev       *C.struct_libevdev
	fd        *os.File
	modifiers *KeyModifiers
}

func (k *KeyboardDevice) Grab() (func() error, error) {
	return grabDevice(k.dev)
}

// Close will gracefully handle closing a keyboard device, freeing memory and
// file descriptors
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

// OpenKeyboardDevice will open a specific keyboard device (from the device path
// passed as a string)
func OpenKeyboardDevice(devPath string) (*KeyboardDevice, error) {
	dev := C.libevdev_new()
	fd, err := os.Open(devPath)
	if err != nil {
		return nil, err
	}
	c_err := C.libevdev_set_fd(dev, C.int(fd.Fd()))
	if c_err > 0 {
		return nil, errors.New("failed to init libevdev")
	}
	return &KeyboardDevice{
		dev:       dev,
		fd:        fd,
		modifiers: NewKeyModifers(),
	}, nil
}

// OpenAllKeyboardDevices will open all currently connected keyboards passing
// them out through a channel for further processing
func OpenAllKeyboardDevices() <-chan *KeyboardDevice {
	kbdChan := make(chan *KeyboardDevice)
	go func() {
		for _, kbdPath := range findAllInputDevices() {
			kbd, err := OpenKeyboardDevice(kbdPath)
			if err != nil {
				log.Error().Err(err).
					Msgf("Unable to open device %s.", kbdPath)
			}
			if kbd.isKeyboard() {
				log.Debug().Caller().
					Msgf("Opening keyboard device %s.", kbdPath)
				kbdChan <- kbd
			} else {
				kbd.Close()
			}
		}
		close(kbdChan)
	}()
	return kbdChan
}

func findAllInputDevices() []string {
	var paths []string
	fileRegexp, _ := regexp.Compile(`event\d+$`)
	log.Debug().Caller().
		Msg("Looking for keyboards...")
	err := filepath.WalkDir(devicePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Caller().Err(err).
				Msgf("could not read %q.", path)
			return err
		}
		if !d.IsDir() {
			if fileRegexp.MatchString(path) {
				paths = append(paths, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Caller().Err(err).
			Msgf("Couldn't traverse device path: %s.", devicePath)
	}
	log.Debug().Caller().
		Msg("Keyboard search finished.")
	return paths
}

// SnoopAllKeyboards will snoop or listen for all key events on all currently connected keyboards.  Keyboards are passed in through a channel, see OpenKeyboardDevices for an example of opening all connected keyboards
func SnoopAllKeyboards(ctx context.Context, kbds <-chan *KeyboardDevice) <-chan KeyEvent {
	keys := make(chan KeyEvent)
	var wg sync.WaitGroup
	for kbd := range kbds {
		log.Debug().Caller().
			Msgf("Tracking keys on device %s.", kbd.fd.Name())
		wg.Add(1)
		go func(ctx context.Context, k *KeyboardDevice, keyCh chan KeyEvent) {
			defer wg.Done()
			kbdSnoop(ctx, k, keyCh)
		}(ctx, kbd, keys)
	}
	go func() {
		<-ctx.Done()
		close(keys)
		wg.Wait()
	}()
	return keys
}

// SnoopKeyboard will snoop or listen for all key events on the given keyboard
// device.
func SnoopKeyboard(ctx context.Context, kbd *KeyboardDevice) <-chan KeyEvent {
	keys := make(chan KeyEvent)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		kbdSnoop(ctx, kbd, keys)
	}()
	go func() {
		<-ctx.Done()
		close(keys)
		wg.Wait()
	}()
	return keys
}

func kbdSnoop(ctx context.Context, kbd *KeyboardDevice, keys chan KeyEvent) {
	norm := C.enum_libevdev_read_flag(C.LIBEVDEV_READ_FLAG_NORMAL)
	for {
		select {
		case <-ctx.Done():
			return
		default:
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
}

// VirtualKeyboardDevice represents a "virtual" (uinput) keyboard device
type VirtualKeyboardDevice struct {
	uidev   *C.struct_libevdev_uinput
	dev     *C.struct_libevdev
	Name    string
	DevNode string
	SysPath string
}

// NewVirtualKeyboard will create a new virtual keyboard device (with the name
// passed in)
func NewVirtualKeyboard(name string) (*VirtualKeyboardDevice, error) {
	if name == "" {
		return nil, errors.New("no name provided")
	}
	var uidev *C.struct_libevdev_uinput

	uid, gid := getUserIds()
	setIDsWithCaps(0, 0, nil)

	dev := C.libevdev_new()
	C.libevdev_set_name(dev, C.CString(name))
	// expose the relevant event types
	C.libevdev_enable_event_type(dev, C.EV_REL)
	C.libevdev_enable_event_type(dev, C.EV_KEY)
	C.libevdev_enable_event_type(dev, C.EV_REP)
	C.libevdev_enable_event_type(dev, C.EV_SYN)
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
	if rv > 0 || uidev == nil {
		return nil, errors.New("failed to create new uinput device")
	}
	log.Debug().Caller().
		Msgf("Virtual keyboard created at %s.",
			C.GoString(C.libevdev_uinput_get_devnode(uidev)))
	time.Sleep(time.Millisecond * 500)

	setIDsWithCaps(uid, gid, []int{getInputGroupGid()})
	if err := cap.NewSet().SetProc(); err != nil {
		return nil, fmt.Errorf("unable to drop privilege: %v", err)
	}

	return &VirtualKeyboardDevice{
		uidev:   uidev,
		dev:     dev,
		Name:    name,
		DevNode: C.GoString(C.libevdev_uinput_get_devnode(uidev)),
		SysPath: C.GoString(C.libevdev_uinput_get_syspath(uidev)),
	}, nil
}

func (u *VirtualKeyboardDevice) sendKeys(done <-chan struct{}, ev ...<-chan *key) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)
	output := func(in <-chan *key) {
		for k := range in {
			select {
			case <-done:
				return
			default:
				rv := C.libevdev_uinput_write_event(u.uidev, C.uint(k.keyType), C.uint(k.keyCode), C.int(k.value))
				if rv < 0 {
					out <- fmt.Errorf("failed send key event type: %v code: %v value %v", k.keyType, k.keyCode, k.value)
				}
				time.Sleep(time.Microsecond)
			}
		}
		wg.Done()
	}

	wg.Add(len(ev))
	for _, c := range ev {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (u *VirtualKeyboardDevice) TypeKey(c int, holdShift bool) error {
	done := make(chan struct{})
	defer close(done)
	if holdShift {
		errc := u.sendKeys(done, keySequence(keyPress(C.KEY_LEFTSHIFT), keySync(), keyPress(c), keySync(), keyRelease(c), keySync(), keyRelease(C.KEY_LEFTSHIFT), keySync()))
		if err := <-errc; err != nil {
			return err
		}
	} else {
		errc := u.sendKeys(done, keySequence(keyPress(c), keySync(), keyRelease(c), keySync()))
		if err := <-errc; err != nil {
			return err
		}
	}
	return nil
}

func (u *VirtualKeyboardDevice) TypeRune(r rune) error {
	if !unicode.In(r, unicode.PrintRanges...) {
		return fmt.Errorf("rune %c (%U) is not a printable character", r, r)
	}
	keyCode, isUpperCase := CodeAndCase(r)
	if keyCode == 0 {
		return fmt.Errorf("rune %c (%U) is not in rune map", r, r)
	} else {
		return u.TypeKey(keyCode, isUpperCase)
	}
}

// TypeSpace is a high level way to "type" a space character (effectively,
// press/release the spacebar)
func (u *VirtualKeyboardDevice) TypeSpace() error {
	return u.TypeKey(C.KEY_SPACE, false)
}

// TypeBackspace allows you to "type" a backspace key and remove a single
// character
func (u *VirtualKeyboardDevice) TypeBackspace() error {
	return u.TypeKey(C.KEY_BACKSPACE, false)
}

// TypeString is a high level function that makes it easy to "type" out a string
// to the virtual keyboard
func (u *VirtualKeyboardDevice) TypeString(str string) error {
	s := strings.NewReader(str)
	for {
		r, _, err := s.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch r {
		case ' ':
			err := u.TypeSpace()
			if err != nil {
				return err
			}
		default:
			err := u.TypeRune(r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Close will gracefully remove a virtual keyboard, freeing memory and file
// descriptors
func (u *VirtualKeyboardDevice) Close() {
	log.Debug().Caller().
		Msg("Closing virtual keyboard device.")
	C.libevdev_uinput_destroy(u.uidev)
	C.libevdev_free(u.dev)
}

// Grab will grab the virtual keyboard which prevents any other clients and the
// kernel from recieving events from it. The returned func can be used to ungrab
// the keyboard, allowing other clients and the kernel to see its events again.
func (u *VirtualKeyboardDevice) Grab() (func() error, error) {
	kbd, err := OpenKeyboardDevice(u.DevNode)
	if err != nil {
		return nil, fmt.Errorf("could not open %s", u.Name)
	}
	return grabDevice(kbd.dev)
}
