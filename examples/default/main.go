package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
)

func main() {
	conf := gpsgen.NewConfig()

	myDevice, err := conf.NewDevice()
	if err != nil {
		panic(err)
	}
	myDevice.OnStateChange = func(state *gpsgen.State, _ []byte) {
		draw.Table(state)
	}

	gen := gpsgen.New()
	gen.Attach(myDevice)
	gen.Run()
	defer gen.Close()

	select {}
}
