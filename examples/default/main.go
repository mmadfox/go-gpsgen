package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/proto"
)

func main() {
	conf := gpsgen.NewConfig()

	myDevice, err := conf.NewDevice()
	if err != nil {
		panic(err)
	}
	myDevice.OnStateChange = func(state *proto.Device) {
		draw.Table(state)
	}

	gen := gpsgen.New()
	gen.Attach(myDevice)
	gen.Run()
	defer gen.Close()

	select {}
}
