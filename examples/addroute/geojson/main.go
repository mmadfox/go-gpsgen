package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/go-gpsgen"
)

func main() {
	gen := gpsgen.New(nil)

	gen.OnPacket(func(data []byte) {
		pck, err := gpsgen.PacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		tracker := pck.Devices[0]
		fmt.Printf("%s -> %f, %f\n",
			tracker.Model,
			tracker.Location.Lon,
			tracker.Location.Lat)
	})

	droneTracker := gpsgen.NewDroneTracker()

	route, err := gpsgen.GeoJSONDecode([]byte(rawRoute))
	if err != nil {
		panic(err)
	}

	droneTracker.AddRoute(route...)

	terminate(func() {
		gen.Close()
	})

	gen.Attach(droneTracker)
	gen.Run()
}

func terminate(fn func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigChan
		fn()
	}()
}

const rawRoute = `
{
	"type": "FeatureCollection",
	"features": [
	  {
		"type": "Feature",
		"properties": {},
		"geometry": {
		  "coordinates": [
			[
			  -0.6848977007719554,
			  23.766026744995628
			],
			[
			  31.65172140776744,
			  -2.9361412054812916
			],
			[
			  67.67893291456264,
			  44.74332748947083
			],
			[
			  64.86705299207989,
			  48.249999904172114
			],
			[
			  28.31261399982,
			  53.247483804819666
			]
		  ],
		  "type": "LineString"
		}
	  }
	]
  }
`
