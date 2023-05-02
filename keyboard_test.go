// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"testing"
)

func TestOpenKeyboardDevice(t *testing.T) {
	testOpenKeyboardDevice(t)
}

func TestOpenKeyboardDevices(t *testing.T) {
	testOpenKeyboardDevices(t)
}

func TestKeyboardDevice_isKeyboard(t *testing.T) {
	testKeyboardDevice_isKeyboard(t)
}

func TestNewVirtualKeyboard(t *testing.T) {
	testNewVirtualKeyboard(t)
}
