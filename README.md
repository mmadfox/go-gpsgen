<h1 align="center">
  <img src="./.github/gpsgen-logo.png" alt="GPS data generator" height="250px">
  <br>
  GPSGen
  <br>
</h1>

GPS data generator based on predefined routes.
Supports GPX and GeoJSON route formats.

This library can be used in testing and debugging applications or devices dependent on GPS/GLONASS/ETC, allowing you to simulate locations for checking their functionality without actual movement.

<hr />

[![Coverage Status](https://coveralls.io/repos/github/mmadfox/go-gpsgen/badge.svg?branch=main)](https://coveralls.io/github/mmadfox/go-gpsgen?branch=main)
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://pkg.go.dev/github.com/mmadfox/go-gpsgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmadfox/go-gpsgen)](https://goreportcard.com/report/github.com/mmadfox/go-gpsgen)
![Actions](https://github.com/mmadfox/go-gpsgen/actions/workflows/cover.yml/badge.svg)

## Table of Contents

- [Installation](#installation)
- [Example](#example)
- Routes
  - [GeoJSON](#geojson)
  - [GPX](#gpx)
  - [Random](#random)
- [Sensors](#sensors)
- [Generated data](#generated-data)

### Installation

```shell
$ go get github.com/mmadfox/go-gpsgen
```

### Example

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen"
)

const (
	numTracksPerRoute = 3
	numTrackers       = 1000
	flushInterval     = 3 * time.Second
)

func main() {
	lon := 37.625616307117696
	lat := 55.75350460378772

	genOpts := gpsgen.NewOptions()
	genOpts.Interval = flushInterval
	gen := gpsgen.New(genOpts)

	// For network transmission
	gen.OnPacket(func(b []byte) {
		// udp.send(b)
	})

	gen.OnError(func(err error) {
		fmt.Println("[ERROR]", err)
	})

	gen.OnNext(func() {
		fmt.Println("tracker state changed successfully")
	})

    // Generate random routes
	for i := 0; i < numTrackers; i++ {
		tracker := gpsgen.NewTracker()
		tracker.SetUserID(uuid.NewString())
		route := gpsgen.RandomRoute(lon, lat, numTracksPerRoute, gpsgen.RouteLevelM)
		tracker.AddRoute(route)
		gen.Attach(tracker)
	}

	terminate(func() {
		gen.Close()
	})

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
```

### Routes

#### GeoJSON

```go
tracker := gpsgen.NewDroneTracker()

route, err := gpsgen.DecodeGeoJSONRoutes(geoJSONBytes)
if err != nil {
	panic(err)
}

tracker.AddRoute(route...)
// ...
```

#### GPX

```go
tracker := gpsgen.NewDroneTracker()

route, err := gpsgen.DecodeGPXRoutes(GPXBytes)
if err != nil {
	panic(err)
}

tracker.AddRoute(route...)
// ...
```

#### Random

```go
tracker := gpsgen.NewDroneTracker()

// lon, lat, numTracks, zoomLevel
lon := 28.31261399982
lat := 53.247483804819666
numTracks := 2

route := gpsgen.RandomRoute(lon, lat, numTracks, gpsgen.RouteLevelXL)

tracker.AddRoute(route...)
// ...
```

### Sensors

A sensor provides a flexible and expandable way to represent and work with sensors. It enables the generation of different values for various tasks, making it suitable for diverse applications and use cases, including sensor data collection, modeling, or analysis.

```go
// types.WithStart: Generation starts from the minimum value of 1 to 10
// types.WithEnd:   Generation ends at the minimum value from 10 to 1
// types.WithRandom: Data generation follows a Bezier curve from 1 to 10

minValue := 1
maxValue := 10
amplitude := 16 // 4 - 512

s1, err := gpsgen.NewSensor("s1", minValue, maxValue, amplitude, types.WithStart|types.WithRandom|types.WithEnd)
if err != nil {
	panic(err)
}

droneTracker.AddSensor(s1)
s2, err := gpsgen.NewSensor("s2", 10, 20, 16, types.WithRandom|types.WithEnd)
if err != nil {
	panic(err)
}
droneTracker.AddSensor(s2)
s3, err := gpsgen.NewSensor("s3", 20, 30, 16, 0)
if err != nil {
	panic(err)
}
droneTracker.AddSensor(s3)
// ...
```

### Generated Data

```text
Device:
    id
    user_id
    tick
    duration
    model
    speed
    distance
    battery (charge, charge_time)
    routes (routes)
    location (lat, lon, elevation, bearing, lat_dms, lon_dms, utm)
    navigator (current_route_index, current_track_index, current_segment_index)
    sensors (id, name, val_x, val_y)
    description
    is_offline
    offline_duration
    color
    time_estimate
Device.Battery:
    charge
    charge_time
Device.Routes:
    routes (Route)
Device.Routes.Route:
    id
    tracks (Track)
    distance
    color
    props
    props_count
Device.Routes.Route.Track:
    distance
    num_segments
    color
    props
    props_count
Device.Sensor:
    id
    name
    val_x
    val_y
Device.Navigator:
    current_route_index
    current_track_index
    current_segment_index
Device.Distance:
    distance
    current_distance
    route_distance
    current_route_distance
    track_distance
    current_track_distance
    segment_distance
    current_segment_distance
Device.Location:
    lat
    lon
    elevation
    bearing
    lat_dms (degrees, minutes, seconds, direction)
    lon_dms (degrees, minutes, seconds, direction)
    utm (central_meridian, easting, northing, long_zone, lat_zone, hemisphere, srid)
Packet:
    devices (Device)
    timestamp
```
