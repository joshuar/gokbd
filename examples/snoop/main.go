// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joshuar/gokbd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	ctx, cancelFunc := context.WithCancel(context.TODO())

	keys := gokbd.SnoopAllKeyboards(ctx, gokbd.OpenAllKeyboardDevices())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancelFunc()
	}()
	log.Info().Msg("Press Ctrl-C to stop snooping...")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Stopping snooping.")
				return
			case k := <-keys:
				if k.Value == 1 && k.TypeName == "EV_KEY" {
					log.Info().Str("eventType", k.TypeName).Str("eventName", k.EventName).
						Int("eventValue", k.Value).Str("rune", string(k.AsRune)).
						Msg("Key pressed.")
				}
				if k.Value == 0 && k.TypeName == "EV_KEY" {
					log.Info().Str("eventType", k.TypeName).Str("eventName", k.EventName).
						Int("eventValue", k.Value).Str("rune", string(k.AsRune)).
						Msg("Key released.")
				}
				if k.Value == 2 && k.TypeName == "EV_KEY" {
					log.Info().Str("eventType", k.TypeName).Str("eventName", k.EventName).
						Int("eventValue", k.Value).Str("rune", string(k.AsRune)).
						Msg("Key held.")
				}
			}
		}
	}()
	wg.Wait()
}
