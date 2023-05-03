// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyModifiers_ToggleAlt(t *testing.T) {
	type fields struct {
		CapsLock bool
		Alt      bool
		Ctrl     bool
		Shift    bool
		Meta     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "toggle test",
			fields: fields{
				Alt: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := &KeyModifiers{
				CapsLock: tt.fields.CapsLock,
				Alt:      tt.fields.Alt,
				Ctrl:     tt.fields.Ctrl,
				Shift:    tt.fields.Shift,
				Meta:     tt.fields.Meta,
			}
			km.ToggleAlt()
			assert.Equal(t, tt.want, km.Alt)
		})
	}
}

func TestKeyModifiers_ToggleShift(t *testing.T) {
	type fields struct {
		CapsLock bool
		Alt      bool
		Ctrl     bool
		Shift    bool
		Meta     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "toggle test",
			fields: fields{
				Alt: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := &KeyModifiers{
				CapsLock: tt.fields.CapsLock,
				Alt:      tt.fields.Alt,
				Ctrl:     tt.fields.Ctrl,
				Shift:    tt.fields.Shift,
				Meta:     tt.fields.Meta,
			}
			km.ToggleShift()
			assert.Equal(t, tt.want, km.Shift)
		})
	}
}

func TestKeyModifiers_ToggleCtrl(t *testing.T) {
	type fields struct {
		CapsLock bool
		Alt      bool
		Ctrl     bool
		Shift    bool
		Meta     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "toggle test",
			fields: fields{
				Ctrl: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := &KeyModifiers{
				CapsLock: tt.fields.CapsLock,
				Alt:      tt.fields.Alt,
				Ctrl:     tt.fields.Ctrl,
				Shift:    tt.fields.Shift,
				Meta:     tt.fields.Meta,
			}
			km.ToggleCtrl()
			assert.Equal(t, tt.want, km.Ctrl)

		})
	}
}

func TestKeyModifiers_ToggleMeta(t *testing.T) {
	type fields struct {
		CapsLock bool
		Alt      bool
		Ctrl     bool
		Shift    bool
		Meta     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "toggle test",
			fields: fields{
				Meta: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := &KeyModifiers{
				CapsLock: tt.fields.CapsLock,
				Alt:      tt.fields.Alt,
				Ctrl:     tt.fields.Ctrl,
				Shift:    tt.fields.Shift,
				Meta:     tt.fields.Meta,
			}
			km.ToggleMeta()
			assert.Equal(t, tt.want, km.Meta)
		})
	}
}

func TestKeyModifiers_ToggleCapsLock(t *testing.T) {
	type fields struct {
		CapsLock bool
		Alt      bool
		Ctrl     bool
		Shift    bool
		Meta     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "toggle test",
			fields: fields{
				CapsLock: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := &KeyModifiers{
				CapsLock: tt.fields.CapsLock,
				Alt:      tt.fields.Alt,
				Ctrl:     tt.fields.Ctrl,
				Shift:    tt.fields.Shift,
				Meta:     tt.fields.Meta,
			}
			km.ToggleCapsLock()
			assert.Equal(t, tt.want, km.CapsLock)
		})
	}
}

func TestNewKeyModifers(t *testing.T) {
	tests := []struct {
		name string
		want *KeyModifiers
	}{
		{
			name: "default test",
			want: &KeyModifiers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKeyModifers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeyModifers() = %v, want %v", got, tt.want)
			}
		})
	}
}
