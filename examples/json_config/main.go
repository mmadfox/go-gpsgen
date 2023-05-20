package main

import (
	"encoding/json"
	"fmt"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/route"
)

const config = `{
	"model":"RT-A32ClN",
	"userId":"12345678",
	"properties": {
		"foo":"bar",
		"xoo":"joo"
	},
	"description":"some description",
	"speed":{
	   "max":2,
	   "min":1,
	   "amplitude":8
	},
	"battery":{
	   "max":90,
	   "min":50
	},
	"elevation":{
	   "max":30,
	   "min":1,
	   "amplitude":8
	},
	"offline":{
	   "min":1,
	   "max":5
	},
	"sensors": [
		{"name": "s2", "min": 1, "max":3, "amplitude": 8},
		{"name": "s3", "min": 5, "max":10, "amplitude": 16}
	]
 }`

func main() {
	data, _ := json.Marshal(gpsgen.DeviceConstraints())
	fmt.Println(string(data))

	conf := gpsgen.NewConfig()
	if err := json.Unmarshal([]byte(config), conf); err != nil {
		panic(err)
	}

	routes, err := route.RoutesForChina()
	if err != nil {
		panic(err)
	}
	conf.Routes = routes

	dev, err := conf.NewDevice()
	if err != nil {
		panic(err)
	}
	dev.OnStateChange = func(state *gpsgen.State, _ []byte) {
		draw.Table(state)
	}

	gen := gpsgen.New()
	gen.Attach(dev)

	gen.Run()
	defer gen.Close()

	select {}
}
