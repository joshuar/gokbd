package main

import (
	"github.com/joshuar/gokbd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	vDev := gokbd.NewVirtualKeyboard("gokbd test")
	vDev.TypeRune('H')
	vDev.TypeRune('e')
	vDev.TypeRune('l')
	vDev.TypeRune('l')
	vDev.TypeRune('o')
	vDev.TypeSpace()
	vDev.TypeString("there!")
	vDev.TypeSpace()
	vDev.TypeString("I made a misteak")
	for i := 0; i <= 2; i++ {
		vDev.TypeBackspace()
	}
	vDev.TypeString("ake")

	vDev.Close()
}
