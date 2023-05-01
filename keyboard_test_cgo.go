// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"C"
	"os"
)
import (
	"reflect"
	"testing"
)

func TestKeyboardDevice_Close(t *testing.T) {
	type fields struct {
		dev       *C.struct_libevdev
		fd        *os.File
		modifiers *KeyModifiers
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeyboardDevice{
				dev:       tt.fields.dev,
				fd:        tt.fields.fd,
				modifiers: tt.fields.modifiers,
			}
			k.Close()
		})
	}
}

func TestKeyboardDevice_isKeyboard(t *testing.T) {
	type fields struct {
		dev       *C.struct_libevdev
		fd        *os.File
		modifiers *KeyModifiers
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KeyboardDevice{
				dev:       tt.fields.dev,
				fd:        tt.fields.fd,
				modifiers: tt.fields.modifiers,
			}
			if got := k.isKeyboard(); got != tt.want {
				t.Errorf("KeyboardDevice.isKeyboard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testOpenKeyboardDevice(t *testing.T) {
	kbds := findAllKeyboards()
	type args struct {
		devPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *KeyboardDevice
		wantErr bool
	}{
		{
			name:    "test with keyboard",
			args:    args{devPath: kbds[0]},
			want:    &KeyboardDevice{},
			wantErr: false,
		},
		{
			name:    "test with no keyboard",
			args:    args{devPath: "/fake/path"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OpenKeyboardDevice(tt.args.devPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenKeyboardDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("OpenKeyboardDevice() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func testOpenKeyboardDevices(t *testing.T) {
	tests := []struct {
		name string
		want <-chan *KeyboardDevice
	}{
		{
			name: "test open",
			want: make(<-chan *KeyboardDevice),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := OpenKeyboardDevices(); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("OpenKeyboardDevices() = %v, want %v", got, tt.want)
			// }
			if got := OpenKeyboardDevices(); got == nil {
				t.Errorf("OpenKeyboardDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnoopAllKeyboards(t *testing.T) {
	type args struct {
		kbds <-chan *KeyboardDevice
	}
	tests := []struct {
		name string
		args args
		want <-chan KeyEvent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnoopAllKeyboards(tt.args.kbds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SnoopAllKeyboards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVirtualKeyboard(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *VirtualKeyboardDevice
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVirtualKeyboard(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVirtualKeyboard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyPress(t *testing.T) {
	type args struct {
		c int
	}
	tests := []struct {
		name string
		args args
		want *key
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyPress(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyPress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyRelease(t *testing.T) {
	type args struct {
		c int
	}
	tests := []struct {
		name string
		args args
		want *key
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyRelease(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keySync(t *testing.T) {
	tests := []struct {
		name string
		want *key
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keySync(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keySync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keySequence(t *testing.T) {
	type args struct {
		keys []*key
	}
	tests := []struct {
		name string
		args args
		want <-chan *key
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keySequence(tt.args.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keySequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualKeyboardDevice_TypeKey(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	type args struct {
		c         int
		holdShift bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.TypeKey(tt.args.c, tt.args.holdShift)
		})
	}
}

func TestVirtualKeyboardDevice_TypeRune(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.TypeRune(tt.args.r)
		})
	}
}

func TestVirtualKeyboardDevice_sendKeys2(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	type args struct {
		keys []*key
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.sendKeys2(tt.args.keys...)
		})
	}
}

func TestVirtualKeyboardDevice_sendKeys(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	type args struct {
		done <-chan struct{}
		ev   []<-chan *key
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			if got := u.sendKeys(tt.args.done, tt.args.ev...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualKeyboardDevice.sendKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualKeyboardDevice_TypeSpace(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.TypeSpace()
		})
	}
}

func TestVirtualKeyboardDevice_TypeBackspace(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.TypeBackspace()
		})
	}
}

func TestVirtualKeyboardDevice_TypeString(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.TypeString(tt.args.str)
		})
	}
}

func TestVirtualKeyboardDevice_Close(t *testing.T) {
	type fields struct {
		uidev *C.struct_libevdev_uinput
		dev   *C.struct_libevdev
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &VirtualKeyboardDevice{
				uidev: tt.fields.uidev,
				dev:   tt.fields.dev,
			}
			u.Close()
		})
	}
}
