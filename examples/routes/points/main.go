package main

import "github.com/mmadfox/go-gpsgen/navigator"

func main() {
	r1, err := navigator.NewRoute([][]navigator.Point{
		{
			{X: 55.748482140161286, Y: 37.62444198526788}, // X - lat, Y - lon
			{X: 55.74863130542925, Y: 37.62442905718493},  // X - lat, Y - lon
		},
	})
	if err != nil {
		panic(err)
	}
	_ = r1
}
