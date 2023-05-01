// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"os"
	"time"

	"github.com/joshuar/gokbd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	vDev := gokbd.NewVirtualKeyboard("gokbd test")
	vDev.TypeRune('H')
	vDev.TypeRune('e')
	vDev.TypeRune('l')
	vDev.TypeRune('l')
	vDev.TypeRune('o')
	vDev.TypeSpace()
	time.Sleep(time.Millisecond * 500)
	vDev.TypeString("there!")
	vDev.TypeSpace()
	vDev.TypeString("I made a misteak")
	time.Sleep(time.Second)
	for i := 0; i <= 2; i++ {
		vDev.TypeBackspace()
	}
	time.Sleep(time.Second)
	vDev.TypeString("ake")
	time.Sleep(time.Second * 2)

	vDev.Close()
}
