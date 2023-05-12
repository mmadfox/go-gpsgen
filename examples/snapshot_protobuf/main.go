package main

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/route"
	"google.golang.org/protobuf/proto"
)

func main() {
	r1, err := route.China3()
	if err != nil {
		panic(err)
	}

	dataCh := make(chan []byte, 1)
	go func() {
		for data := range dataCh {
			snap := new(pb.Device)
			if err := proto.Unmarshal(data, snap); err != nil {
				panic(err)
			}
			fmt.Printf("model=%s, lat=%f, lon=%f, el=%f, curDist=%f\n",
				snap.Model,
				snap.Latitude,
				snap.Longitude,
				snap.Elevation,
				snap.CurrentDistance,
			)
		}
	}()

	tracker1, err := gpsgen.Tracker("example", nil, r1)
	if err != nil {
		panic(err)
	}
	tracker1.OnStateChange = func(_ *gpsgen.State, snapshot []byte) {
		dataCh <- snapshot
	}

	gen := gpsgen.New()
	gen.Attach(tracker1)

	gen.Run()

	select {}
}
