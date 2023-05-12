package main

import "github.com/mmadfox/go-gpsgen/route"

func main() {
	myRoute, err := route.Generate()
	if err != nil {
		panic(err)
	}
	_ = myRoute

	myRouteR1, err := route.GenerateFor(route.Russia)
	if err != nil {
		panic(err)
	}
	_ = myRouteR1
}
