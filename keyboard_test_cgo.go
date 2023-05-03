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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// handy code for testing channel return from:
// https://www.sidorenko.io/post/2019/01/testing-of-functions-with-channels-in-go/
func getValues[V any](c <-chan V) []V {
	var r []V
	for i := range c {
		r = append(r, i)
	}
	return r
}

func testKeyboardDevice_isKeyboard(t *testing.T) {
	virtualKbd, err := NewVirtualKeyboard("gokbdtest")
	assert.Nil(t, err)
	kbds := OpenAllKeyboardDevices()
	realKbd := <-kbds
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
		{
			name: "is a keyboard",
			fields: fields{
				dev: realKbd.dev,
			},
			want: true,
		},
		{
			name: "not a keyboard",
			fields: fields{
				dev: virtualKbd.dev,
			},
			want: false,
		},
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
	kbds := findAllInputDevices()
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
			if got := OpenAllKeyboardDevices(); got == nil {
				t.Errorf("OpenKeyboardDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testNewVirtualKeyboard(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *VirtualKeyboardDevice
		wantErr bool
	}{
		{
			name: "test successful creation",
			args: args{name: "gokdb-successful-test"},
		},
		{
			name:    "test empty string",
			args:    args{name: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVirtualKeyboard(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVirtualKeyboard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewVirtualKeyboard() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func testVirtualKeyboardDevice_TypeKey(t *testing.T) {
	type args struct {
		c         int
		holdShift bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test type key (no shift)",
			args:    args{c: 30, holdShift: false},
			wantErr: false,
		},
		{
			name:    "test type key (shift)",
			args:    args{c: 30, holdShift: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVirtualKbd, err := NewVirtualKeyboard(tt.name)
			assert.Nil(t, err)
			if err := testVirtualKbd.TypeKey(tt.args.c, tt.args.holdShift); (err != nil) != tt.wantErr {
				t.Errorf("VirtualKeyboardDevice.TypeKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			testVirtualKbd.Close()
		})
	}
}

func testVirtualKeyboardDevice_TypeRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test type lowercase",
			args:    args{r: 'a'},
			wantErr: false,
		},
		{
			name:    "test type uppercase",
			args:    args{r: 'A'},
			wantErr: false,
		},
		{
			name:    "test not in rune map",
			args:    args{r: '🦍'},
			wantErr: true,
		},
		{
			name:    "test space",
			args:    args{r: ' '},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewVirtualKeyboard(tt.name)
			assert.Nil(t, err)
			if err := u.TypeRune(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("VirtualKeyboardDevice.TypeRune() error = %v, wantErr %v", err, tt.wantErr)
			}
			u.Close()
		})
	}
}

func testVirtualKeyboardDevice_TypeString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test string no space",
			args:    args{str: "string"},
			wantErr: false,
		},
		{
			name:    "test string space",
			args:    args{str: "string space"},
			wantErr: false,
		},
		{
			name:    "test invalid rune in string",
			args:    args{str: "str🦍ng"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewVirtualKeyboard(tt.name)
			assert.Nil(t, err)
			if err := u.TypeString(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("VirtualKeyboardDevice.TypeString() error = %v, wantErr %v", err, tt.wantErr)
			}
			u.Close()
		})
	}
}
