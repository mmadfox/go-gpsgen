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
