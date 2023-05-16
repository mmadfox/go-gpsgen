<h1 align="center">
  <img src="./.github/gpsgen-logo.png" alt="GPS data generator" height="250px">
  <br>
  GPSGen
  <br>
</h1>

Language: [En](./README.md)

Генератор GPS данных на основе заданных маршрутов.

Библиотека может использоваться в тестировании и отладке приложений или устройств, зависящих от GPS, позволяя создавать симулированные местоположения для проверки их функциональности без необходимости фактического перемещения.

## Оглавление
+ [Примеры](#примеры)
+ [Установка](#установка)
+ [Генерируемые данные](#генерируемые-данные)
+ [Ограничения и базовые единицы](#ограничения-и-базовые-единицы)
+ [Настройки](#настройки)
+ [Маршруты](#маршруты)
  - [GeoJSON](#geojson)
  - [GPX](#gpx)
  - [Рандомные](#рандомный)
  - [Статичные](#cтатичный)
  - [Стандартный](#стандартный)
+ [Пресеты устройств](#пресеты-устройств)  


### Примеры
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
### Установка
```shell
$ go get github.com/mmadfox/go-gpsgen
```

### Генерируемые данные
| Name              | Data                                        |
|-------------------|---------------------------------------------|
| Device            | model, description, properties              |
| Distance (meters) | total, current                              |
| Speed (m/s)       | current speed                               |
| Navigator         | lat, lon, alt, bearing, DMSLat, DMSLon, UTM |
| Sensor            | value_x, value_y, name                      |
| Status            | online, offline                             |
| User              | custom id, tick                             |

### Ограничения и базовые единицы
| Option        | Constraint                          | Unit             |
|---------------|-------------------------------------|------------------|
| WithSpeed     | min=0, max=1000, amplitude=4..512   | meter per second |
| WithBattery   | min=0, max=100                      | percent          |
| WithModel     | min=1, max=64                       |                  |
| WithElevation | min=0, max=100000, amplitude=4..512 | meters           |
| WithOffline   | min=0, max=300                      | seconds          |
| WithSensors   | amplitude=4..512                    | any              |

### Настройки
| Option                | Description                                                                                     |
|-----------------------|-------------------------------------------------------------------------------------------------|
| WithModel             | Название модели устройства (TxNo12-Oi7)                                                         |
| WithUserID            | Данные пользователя (ID, Number, etc)                                                           |
| WithRoutes, WithRoute | Маршрут для генератора                                                                          |
| WithSpeed             | Диапазон скорости min..max в метрах в секунду                                                   |
| WithBattery           | Диапазон заряда аккумулятора 0..100 в процентах                                                 |
| WithSensors           | Дополнительные датчики, которые можно добавить к устройству. Sensor{name, min, max, amplitude}  |
| WithElevation         | Диапазон высоты min..max в метрах                                                               |
| WithOffline           | Между переключением треков или маршрутов устройство будет уходить в offline min..max в секундах |
| WithProps             | Характеристики устройства map[string]string                                                     |
| WithDescritpion       | Короткое описание устройства                                                                    |

### Маршруты
Маршрут для GPS генератора - это заданный путь или последовательность точек на карте, которые определяют планируемое или фиксированное перемещение геотрекера. 
#### GeoJSON
```go
package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/geojson"
)

func main() {
	geoJSONRoute := `{"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"coordinates":[[37.73185507353088,55.59010742814678],[37.740460827380986,55.585012323942266]],"type":"LineString"}}]}`
	routes, err := geojson.Decode([]byte(geoJSONRoute))
	if err != nil {
		panic(err)
	}
	myTracker, err := gpsgen.Tracker("Tracker78-1", nil, routes...)
	if err != nil {
		panic(err)
	}
	_ = myTracker
}
```
#### GPX
```go 
package main

import "github.com/mmadfox/go-gpsgen/gpx"

func main() {
	gpxRoute := `<?xml version="1.0" encoding="UTF-8"?><gpx xmlns="http://www.topografix.com/GPX/1/1"
   version="1.1"
   creator="Wikipedia" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd"><time>2011-09-22T18:56:51Z</time><metadata><name>Name</name><desc>Description</desc><author><name>Autor</name></author></metadata><rte><rtept lat="55.74966429698134" lon="37.624339525581576"/><rtept lat="55.748482140161286" lon="37.62444198526788"/></rte></gpx>`

	routes, err := gpx.Decode([]byte(gpxRoute))
	if err != nil {
		panic(err)
	}
	_ = routes
}
```
#### Рандомный
Поддерживаются страны: ```Angola, SouthArabia, Turkey, Russia, France, Spain, China```
```go 
package main

import "github.com/mmadfox/go-gpsgen/route"

func main() {
	myRoute, err := route.Generate()
	if err != nil {
		panic(err)
	}
	_ = myRoute

    // by specified country 
	myRoute1, err := route.GenerateFor(route.Russia)
	if err != nil {
		panic(err)
	}
	_ = myRouteR1
}
```
#### Cтатичный
Поддерживаются страны: ```Russia1..5, France1..5, Spain1..5, China1..5``` 
```go
package main

import "github.com/mmadfox/go-gpsgen/route"

func main() {
	_, _ = route.China1()
	_, _ = route.China2()
	_, _ = route.China3()
	_, _ = route.China4()
	_, _ = route.China5()

	_, _ = route.France1()
	_, _ = route.France2()
	_, _ = route.France3()
	_, _ = route.France4()
	_, _ = route.France5()

	_, _ = route.Spain1()
	_, _ = route.Spain2()
	_, _ = route.Spain3()
	_, _ = route.Spain4()
	_, _ = route.Spain5()

	_, _ = route.Russia1()
	_, _ = route.Russia2()
	_, _ = route.Russia3()
	_, _ = route.Russia4()
	_, _ = route.Russia5()
}
```
#### Стандартный
```go
package main

import "github.com/mmadfox/go-gpsgen/navigator"

func main() {
	r1, err := navigator.NewRoute([][]navigator.Point{
		{ // Track-1
			{X: 55.748482140161286, Y: 37.62444198526788}, // X - lat, Y - lon
			{X: 55.74863130542925, Y: 37.62442905718493},  // X - lat, Y - lon
            // ...
		},
        { // Track-2
            // ...
        }
	})
	if err != nil {
		panic(err)
	}
	_ = r1
}
```
### Пресеты устройств
```go 
// Drone 
gpsgen.Drone(model string, props gpsgen.Properties, route ...*navigator.Route) (*gpsgen.Device, error)
gpsgen.DroneWithSensors( model string, routes []*navigator.Route, props gpsgen.Properties, sensors ...gpsgen.Sensor) (*gpsgen.Device, error)

// Tracker
gpsgen.Tracker( model string, props gpsgen.Properties, route ...*navigator.Route) (*gpsgen.Device, error)
gpsgen.TrackerWithSensors( model string, routes []*navigator.Route, props gpsgen.Properties, sensors ...gpsgen.Sensor) (*gpsgen.Device, error)
``` 
### Пример amplitude option
На рисунках представлены рандомные значения от 0 до 120 и контрольные точки 4,8,16,32,64,128,256,512

<img src="./.github/amplitudes/4.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/8.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/16.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/32.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/64.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/128.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/256.png" alt="GPS data generator" height="150px">
<img src="./.github/amplitudes/512.png" alt="GPS data generator" height="150px">
