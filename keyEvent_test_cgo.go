// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"C"
)
import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testNewKeyEvent(t *testing.T) {
	// This is a bit janky, but no easy way to generate the raw key event a
	// keyboard produces without actually having a keyboard generate the key...
	v, err := NewVirtualKeyboard("gokbd test")
	assert.Nil(t, err)
	k, err := OpenKeyboardDevice(v.DevNode)
	assert.Nil(t, err)
	_, err = k.Grab()
	assert.Nil(t, err)
	keyChan := SnoopKeyboard(context.TODO(), k)
	var wantKey KeyEvent
	go func() {
		wantKey = <-keyChan
	}()
	err = v.TypeRune('a')
	assert.Nil(t, err)
	v.Close()

	type args struct {
		ev C.struct_input_event
	}
	tests := []struct {
		name string
		args args
		want *KeyEvent
	}{
		{
			name: "test key event",
			args: args{ev: wantKey.eventRaw},
			want: &wantKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKeyEvent(tt.args.ev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeyEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testKeyEvent_updateRune(t *testing.T) {
	// This is a bit janky, but no easy way to generate the raw key event a
	// keyboard produces without actually having a keyboard generate the key...
	v, err := NewVirtualKeyboard("gokbd test")
	assert.Nil(t, err)
	k, err := OpenKeyboardDevice(v.DevNode)
	assert.Nil(t, err)
	_, err = k.Grab()
	assert.Nil(t, err)
	keyChan := SnoopKeyboard(context.TODO(), k)
	var wantKey KeyEvent
	go func() {
		wantKey = <-keyChan
	}()
	err = v.TypeRune('a')
	assert.Nil(t, err)
	v.Close()

	type fields struct {
		key KeyEvent
	}
	type args struct {
		modifiers *KeyModifiers
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rune
	}{
		{
			name: "test shift",
			fields: fields{
				key: wantKey,
			},
			args: args{modifiers: &KeyModifiers{
				Shift: true,
			}},
			want: 'A',
		},
		{
			name: "test capslock",
			fields: fields{
				key: wantKey,
			},
			args: args{modifiers: &KeyModifiers{
				CapsLock: true,
			}},
			want: 'A',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.key
			k.updateRune(tt.args.modifiers)
			assert.Equal(t, tt.want, k.AsRune)
		})
	}
}

func testKeyEvent_IsKeyPress(t *testing.T) {
	type fields struct {
		Value    int
		TypeName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test valid",
			fields: fields{Value: 1, TypeName: "EV_KEY"},
			want:   true,
		},
		{
			name:   "test invalid",
			fields: fields{Value: 0, TypeName: "EV_KEY"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kev := &KeyEvent{
				Value:    tt.fields.Value,
				TypeName: tt.fields.TypeName,
			}
			if got := kev.IsKeyPress(); got != tt.want {
				t.Errorf("KeyEvent.IsKeyPress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testKeyEvent_IsKeyRelease(t *testing.T) {
	type fields struct {
		Value    int
		TypeName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test valid",
			fields: fields{Value: 0, TypeName: "EV_KEY"},
			want:   true,
		},
		{
			name:   "test invalid",
			fields: fields{Value: 1, TypeName: "EV_KEY"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kev := &KeyEvent{
				Value:    tt.fields.Value,
				TypeName: tt.fields.TypeName,
			}
			if got := kev.IsKeyRelease(); got != tt.want {
				t.Errorf("KeyEvent.IsKeyRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testKeyEvent_IsBackspace(t *testing.T) {
	type fields struct {
		EventName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test valid",
			fields: fields{EventName: "KEY_BACKSPACE"},
			want:   true,
		},
		{
			name:   "test invalid",
			fields: fields{EventName: "EV_KEY"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kev := &KeyEvent{
				EventName: tt.fields.EventName,
			}
			if got := kev.IsBackspace(); got != tt.want {
				t.Errorf("KeyEvent.IsBackspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testKeyEvent_IsModifier(t *testing.T) {
	type fields struct {
		EventName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test valid",
			fields: fields{EventName: "KEY_LEFTSHIFT"},
			want:   true,
		},
		{
			name:   "test invalid",
			fields: fields{EventName: "EV_KEY"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kev := &KeyEvent{
				EventName: tt.fields.EventName,
			}
			if got := kev.IsModifier(); got != tt.want {
				t.Errorf("KeyEvent.IsModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
