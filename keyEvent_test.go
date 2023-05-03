// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import "testing"

// * "go test" does not support cgo, so need to "fake" the test functions.
// * See keyEvent_test_cgo.go for actual test code.

func TestNewKeyEvent(t *testing.T) {
	testNewKeyEvent(t)
}

func TestKeyEvent_updateRune(t *testing.T) {
	testKeyEvent_updateRune(t)
}

func TestKeyEvent_IsKeyPress(t *testing.T) {
	testKeyEvent_IsKeyPress(t)
}

func TestKeyEvent_IsKeyRelease(t *testing.T) {
	testKeyEvent_IsKeyRelease(t)
}

func TestKeyEvent_IsBackspace(t *testing.T) {
	testKeyEvent_IsBackspace(t)
}

func TestKeyEvent_IsModifier(t *testing.T) {
	testKeyEvent_IsModifier(t)
}
