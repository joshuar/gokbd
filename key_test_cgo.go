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
	"reflect"
	"testing"
)

func test_keyPress(t *testing.T) {
	aKey := &key{
		keyType: C.EV_KEY,
		keyCode: 30,
		value:   1,
	}
	type args struct {
		c int
	}
	tests := []struct {
		name string
		args args
		want *key
	}{
		{
			name: "test press",
			args: args{c: 30},
			want: aKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyPress(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyPress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func test_keyRelease(t *testing.T) {
	aKey := &key{
		keyType: C.EV_KEY,
		keyCode: 30,
		value:   0,
	}
	type args struct {
		c int
	}
	tests := []struct {
		name string
		args args
		want *key
	}{
		{
			name: "test release",
			args: args{c: 30},
			want: aKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyRelease(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func test_keySync(t *testing.T) {
	k := &key{
		keyType: C.EV_SYN,
		keyCode: C.SYN_REPORT,
		value:   0,
	}
	tests := []struct {
		name string
		want *key
	}{
		{
			name: "test sync",
			want: k,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keySync(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keySync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func test_keySequence(t *testing.T) {
	keyList := []*key{
		keyPress(30),
		keySync(),
		keyRelease(30),
		keySync(),
	}
	type args struct {
		keys []*key
	}
	tests := []struct {
		name string
		args args
		want []*key
	}{
		{
			name: "test key sequence",
			args: args{keys: keyList},
			want: keyList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValues(keySequence(tt.args.keys...)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keySequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
