// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"os"

	gokbd "github.com/joshuar/gokbd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	keys := gokbd.SnoopAllKeyboards(gokbd.OpenKeyboardDevices())
	for k := range keys {
		if k.Value == 1 && k.TypeName == "EV_KEY" {
			log.Debug().Msgf("Key pressed: %s %s %d %c\n", k.TypeName, k.EventName, k.Value, k.AsRune)
		}
		if k.Value == 0 && k.TypeName == "EV_KEY" {
			log.Debug().Msgf("Key released: %s %s %d\n", k.TypeName, k.EventName, k.Value)
		}
		if k.Value == 2 && k.TypeName == "EV_KEY" {
			log.Debug().Msgf("Key held: %s %s %d %c\n", k.TypeName, k.EventName, k.Value, k.AsRune)
		}
	}
}
