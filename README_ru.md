<h1 align="center">
  <img src="./.github/gpsgen-logo.png" alt="GPS data generator" height="250px">
  <br>
  GPSGen
  <br>
</h1>

Language: [En](./README.md)

Генератор GPS данных на основе заданных маршрутов.
Поддерживаются маршруты GPX, GeoJSON формата.

Библиотека может использоваться в тестировании и отладке приложений или устройств, зависящих от GPS/GLONASS/ETC, позволяя создавать симулированные местоположения для проверки их функциональности без необходимости фактического перемещения.

<hr />

[![Coverage Status](https://coveralls.io/repos/github/mmadfox/go-gpsgen/badge.svg?branch=main)](https://coveralls.io/github/mmadfox/go-gpsgen?branch=main)
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://pkg.go.dev/github.com/mmadfox/go-gpsgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmadfox/go-gpsgen)](https://goreportcard.com/report/github.com/mmadfox/go-gpsgen)
![Actions](https://github.com/mmadfox/go-gpsgen/actions/workflows/cover.yml/badge.svg)

## Оглавление

- [Установка](#установка)
- [Пример](#пример)
- Маршруты
  - [GeoJSON](#geojson)
  - [GPX](#gpx)
  - [Random](#random)
- [Сенсоры](#сенсоры)
- [Генерируемые данные](#генерируемые-данные)

### Установка

```shell
$ go get github.com/mmadfox/go-gpsgen
```

### Пример

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

	// Для передечи пакета по сети
	gen.OnPacket(func(b []byte) {
		// udp.send(b)
	})

	gen.OnError(func(err error) {
		fmt.Println("[ERROR]", err)
	})

	gen.OnNext(func() {
		fmt.Println("tracker state changed successfully")
	})

    // Генерируем рандомные маршруты
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

### Маршруты

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

### Сенсоры

Сенсор обеспечивает гибкий и расширяемый способ представления датчиков и работы с ними.
Он позволяет генерировать разные значения для разных задач, что делает его подходящим
для различных приложений и вариантов использования, включающих сбор данных датчиков, моделирование или анализ.

```go
// types.WithStart: Генерация начинается с минимального значения от 1 до 10
// types.WithEnd:   Генерация заканчивается минимальным значением от 10 до 1
// types.WithRandom: Генерация данных по кривой безье 1 - 10

minValue := 1
maxValue := 10
amplitude := 16 // 4 - 512

droneTracker.AddSensor("s1", minValue, maxValue, amplitude, types.WithStart|types.WithRandom|types.WithEnd)
droneTracker.AddSensor("s2", 10, 20, 16, types.WithRandom|types.WithEnd)
droneTracker.AddSensor("s3", 20, 30, 16, 0)
```

### Генерируемые данные

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
