package navigator

import (
	"fmt"
	"math"

	"github.com/icholy/utm"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`

	Bearing         float64 `json:"bearing"`
	CurrentDistance float64 `json:"currentDistance"`
	TotalDistance   float64 `json:"routeDistance"`

	LatDMS DMS `json:"latDMS"`
	LonDMS DMS `json:"lonDMS"`
	UTM    UTM `json:"utm"`
}

type DMS struct {
	Degrees   int     `json:"degrees"`
	Minutes   int     `json:"minutes"`
	Seconds   float64 `json:"seconds"`
	Direction string  `json:"direction"`
}

func (d DMS) String() string {
	return fmt.Sprintf(`%d°%d'%f" %s`,
		d.Degrees, d.Minutes, d.Seconds, d.Direction)
}

func ToDMS(lat, lon float64) (latDMS DMS, lonDMS DMS) {
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

	latDMS = DMS{Degrees: latitude, Minutes: latitudeMinutes, Seconds: latitudeSeconds, Direction: latDir}
	lonDMS = DMS{Degrees: longitude, Minutes: longitudeMinutes, Seconds: longitudeSeconds, Direction: lonDir}

	return

}

type UTM struct {
	CentralMeridian float64 `json:"centralMeridian"`
	Easting         float64 `json:"easting"`
	Northing        float64 `json:"northing"`
	LongZone        int     `json:"longZone"`
	LatZone         string  `json:"latZone"`
	Hemisphere      string  `json:"hemisphere"`
	SRID            int     `json:"SRIDCode"`
}

func (u UTM) String() string {
	return fmt.Sprintf("UTM{LongZone: %d, LatZone: %s, Hemisphere: %s, Easting: %f, Northing: %f}",
		u.LongZone, u.LatZone, u.Hemisphere, u.Easting, u.Northing)
}

func ToUTM(lat, lon float64) UTM {
	e, n, z := utm.ToUTM(lat, lon)
	hemisphere := "S"
	if z.North {
		hemisphere = "N"
	}
	return UTM{
		Easting:    e,
		Northing:   n,
		LongZone:   z.Number,
		LatZone:    string(z.Letter),
		Hemisphere: hemisphere,
		SRID:       z.SRID(),
	}
}
