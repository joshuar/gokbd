// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

import (
	"testing"
)

// ?INFO "go test" does not support cgo, so need to "fake" the test functions.
// ? See key_test_cgo.go for actual test code.

func Test_keyPress(t *testing.T) {
	test_keyPress(t)
}

func Test_keyRelease(t *testing.T) {
	test_keyPress(t)
}

func Test_keySync(t *testing.T) {
	test_keyPress(t)
}

func Test_keySequence(t *testing.T) {
	test_keySequence(t)
}
