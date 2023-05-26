package navigator

import (
	"fmt"
	"math"

	"github.com/icholy/utm"
	"github.com/mmadfox/go-gpsgen/proto"
)

func SetDMS(lat, lon float64, latDMS *proto.DMS, lonDMS *proto.DMS) {
	var latDir, lonDir string
	if lat > 0 {
		latDir = "N"
	} else {
		latDir = "S"
	}

	if lon > 0 {
		lonDir = "E"
	} else {
		lonDir = "W"
	}

	lat = math.Abs(lat)
	lon = math.Abs(lon)

	latitude := int(lat)
	latitudeMinutes := int((lat - float64(latitude)) * 60)
	latitudeSeconds := (lat - float64(latitude) - float64(latitudeMinutes)/60) * 3600

	longitude := int(lon)
	longitudeMinutes := int((lon - float64(longitude)) * 60)
	longitudeSeconds := (lon - float64(longitude) - float64(longitudeMinutes)/60) * 3600

	latDMS.Degrees = int64(latitude)
	latDMS.Minutes = int64(latitudeMinutes)
	latDMS.Seconds = latitudeSeconds
	latDMS.Direction = latDir

	lonDMS.Degrees = int64(longitude)
	lonDMS.Minutes = int64(longitudeMinutes)
	lonDMS.Seconds = longitudeSeconds
	lonDMS.Direction = lonDir
}

func SetUTM(lat, lon float64, u *proto.UTM) {
	e, n, z := utm.ToUTM(lat, lon)
	hemisphere := "S"
	if z.North {
		hemisphere = "N"
	}
	u.CentralMeridian = z.CentralMeridian()
	u.Easting = e
	u.Northing = n
	u.LongZone = int64(z.Number)
	u.LatZone = string(z.Letter)
	u.Hemisphere = hemisphere
	u.Srid = int64(z.SRID())
}

func FormatUTM(u *proto.UTM) string {
	return fmt.Sprintf("UTMEasting:%.4f\nUTMNorthing:%.4f\nUTMZone:%d%s\n",
		u.Easting, u.Northing, u.LongZone, u.LatZone)
}

func FormatDMS(d *proto.DMS) string {
	return fmt.Sprintf("%d°%d'%f\"%s",
		d.Degrees, d.Minutes, d.Seconds, d.Direction)
}
