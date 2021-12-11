package main

import (
	gokbd "github.com/joshuar/gokbd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	keys := gokbd.SnoopAllKeyboards(gokbd.OpenKeyboardDevices())
	for k := range keys {
		if k.Value == 1 && k.TypeName == "EV_KEY" {
			log.Infof("Key pressed: %s %s %d %c\n", k.TypeName, k.EventName, k.Value, k.AsRune)
		}
		if k.Value == 0 && k.TypeName == "EV_KEY" {
			log.Infof("Key released: %s %s %d\n", k.TypeName, k.EventName, k.Value)
		}
		if k.Value == 2 && k.TypeName == "EV_KEY" {
			log.Infof("Key held: %s %s %d %c\n", k.TypeName, k.EventName, k.Value, k.AsRune)
		}
	}
}
