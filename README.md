<h1 align="center">
  <img src="./.github/gpsgen-logo.png" alt="GPS data generator" height="250px">
  <br>
  GPSGen
  <br>
</h1>

[RU](./README_ru.md)

GPS data generator based on given routes.

The library can be used in testing and debugging GPS dependent applications or devices, allowing you to create simulated locations to test their functionality without having to actually move.

## Table of contents
+ [Examples](#examples)
+ [Install](#installation)
+ [Generated data](#generated-data)
+ [Limits and base units](#limits-and-base-units)
+ [Settings](#settings)
+ [Routes](#routes)
   - [GeoJSON](#geojson)
   - [GPX](#gpx)
   - [Random](#random)
   - [Static](#static)
   - [Standard](#standard)

### Examples
[Examples](./examples/)

```go
package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/route"
)

func main() {
	myRoute, err := route.Russia2()
	if err != nil {
		panic(err)
	}

	myDevice, err := gpsgen.NewDevice(
		gpsgen.WithModel("myModel"),
		gpsgen.WithDescritpion("some description"),
		gpsgen.WithElevation(1, 3, 4),
		gpsgen.WithSpeed(1, 7, 64),
		gpsgen.WithOffline(1, 120),
		gpsgen.WithBattery(0, 100),
		gpsgen.WithProps(gpsgen.Properties{
			"foo": "foo",
			"bar": "bar",
		}),
		gpsgen.WithUserID("12345678"),
		gpsgen.WithSensors(
			gpsgen.Sensor{Name: "s1", Min: 1, Max: 10, Amplitude: 8},
			gpsgen.Sensor{Name: "s2", Min: 3, Max: 10, Amplitude: 128}),
		gpsgen.WithRoute(myRoute),
	)
	if err != nil {
		panic(err)
	}

	myDevice.OnStateChange = func(state *gpsgen.State, snapshot []byte) {
		draw.Table(state)
	}

	gen := gpsgen.New()
	gen.Attach(myDevice)

	gen.Run()
	defer gen.Close()

	select {}
}
```

### Install
```shell
$ go get github.com/mmadfox/go-gpsgen
```

### Generated data
| Name              | Data                                        |
|-------------------|---------------------------------------------|
| Device            | model, description, properties              |
| Distance (meters) | total, current                              |
| Speed (m/s)       | current speed                               |
| Navigator         | lat, lon, alt, bearing, DMSLat, DMSLon, UTM |
| Sensor            | value_x, value_y, name                      |
| Status            | online, offline                             |
| User              | custom id, tick                             |

### Limits and base units
| Option        | Constraint                          | Unit             |
|---------------|-------------------------------------|------------------|
| WithSpeed     | min=0, max=1000, amplitude=4..512   | meter per second |
| WithBattery   | min=0, max=100                      | percent          |
| WithModel     | min=1, max=64                       |                  |
| WithElevation | min=0, max=100000, amplitude=4..512 | meters           |
| WithOffline   | min=0, max=300                      | seconds          |
| WithSensors   | amplitude=4..512                    | any              |
