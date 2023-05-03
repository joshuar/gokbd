// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"testing"
)

// ?INFO "go test" does not support cgo, so need to "fake" the test functions.
// ? See keyboard_test_cgo.go for actual test code.

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

func TestVirtualKeyboardDevice_TypeKey(t *testing.T) {
	testVirtualKeyboardDevice_TypeKey(t)
}

func TestVirtualKeyboardDevice_TypeRune(t *testing.T) {
	testVirtualKeyboardDevice_TypeRune(t)
}

func TestVirtualKeyboardDevice_TypeString(t *testing.T) {
	testVirtualKeyboardDevice_TypeString(t)
}
