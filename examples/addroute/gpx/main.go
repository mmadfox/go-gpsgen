package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/go-gpsgen"
)

func main() {
	gen := gpsgen.New(nil)

	gen.OnPacket(func(data []byte) {
		pck, err := gpsgen.PacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		tracker := pck.Devices[0]
		fmt.Printf("%s -> %f, %f\n",
			tracker.Model,
			tracker.Location.Lon,
			tracker.Location.Lat)
	})

	droneTracker := gpsgen.NewDroneTracker()

	route, err := gpsgen.GPXDecode([]byte(rawRoute))
	if err != nil {
		panic(err)
	}

	droneTracker.AddRoute(route...)

	terminate(func() {
		gen.Close()
	})

	gen.Attach(droneTracker)
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

const rawRoute = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.topografix.com/GPX/1/1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.topografix.com/GPX/gpx_style/0/2 http://www.topografix.com/GPX/gpx_style/0/2/gpx_style.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:gpx_style="http://www.topografix.com/GPX/gpx_style/0/2" version="1.1" creator="https://gpx.studio">
<metadata>
    <name>new</name>
    <author>
        <name>gpx.studio</name>
        <link href="https://gpx.studio"></link>
    </author>
</metadata>
<trk>
    <name>new</name>
    <type>Cycling</type>
    <trkseg>
    <trkpt lat="53.830488750300155" lon="35.37857168666637">
        <ele>222.3</ele>
    </trkpt>
    <trkpt lat="54.753350433310544" lon="39.90427062969353">
        <ele>0.0</ele>
    </trkpt>
    <trkpt lat="54.52457211739558" lon="39.90427062969353">
        <ele>0.0</ele>
    </trkpt>
    <trkpt lat="52.19226382198787" lon="43.8148260270665">
        <ele>0.0</ele>
    </trkpt>
    <trkpt lat="54.03735817275171" lon="47.41780965161241">
        <ele>0.0</ele>
    </trkpt>
    <trkpt lat="52.727525690763485" lon="51.67987564650206">
        <ele>0.0</ele>
    </trkpt>
    <trkpt lat="56.09915562942588" lon="56.11769694015005">
        <ele>0.0</ele>
    </trkpt>
    </trkseg>
</trk>
</gpx>
`
