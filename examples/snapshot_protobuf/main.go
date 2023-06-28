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
			fmt.Printf("id=%s, model=%s, lat=%f, lon=%f, el=%f, curDist=%f\n",
				snap.Id,
				snap.Model,
				snap.Location.Lat,
				snap.Location.Lon,
				snap.Location.Alt,
				snap.Location.CurrentDistance,
			)
		}
	}()

	tracker1, err := gpsgen.Tracker("example", nil, r1)
	if err != nil {
		panic(err)
	}
	tracker1.OnStateChangeBytes = func(state []byte) {
		dataCh <- state
	}

	gen := gpsgen.New()
	gen.Attach(tracker1)

	gen.Run()

	select {}
}
