package gpsgen

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen/geojson"
	"github.com/mmadfox/go-gpsgen/gpx"
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/types"
	"google.golang.org/protobuf/proto"
)

// EncodeTracker encodes a Device into a binary format.
func EncodeTracker(t *Device) ([]byte, error) {
	if t == nil {
		return []byte{}, nil
	}
	t.Update()
	return t.MarshalBinary()
}

// DecodeTracker decodes binary data into a Device.
func DecodeTracker(data []byte) (*Device, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gpsgen: no tracker data to decode")
	}
	t := new(Device)
	if err := t.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	t.Update()
	return t, nil
}

// EncodeSensors encodes a slice of Sensor types into binary data.
func EncodeSensors(sensors []*types.Sensor) ([]byte, error) {
	if len(sensors) == 0 {
		return nil, fmt.Errorf("gpsgen: no sensors to encode")
	}
	pbsensors := &pb.Snapshot_Sensors{
		Sensors: make([]*pb.Snapshot_SensorType, len(sensors)),
	}
	for i := 0; i < len(sensors); i++ {
		pbsensors.Sensors[i] = sensors[i].Snapshot()
	}
	return proto.Marshal(pbsensors)
}

// DecodeSensors decodes binary data into a slice of Sensor types.
func DecodeSensors(data []byte) ([]*types.Sensor, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gpsgen: no data to decode sensors")
	}
	pbsensors := new(pb.Snapshot_Sensors)
	if err := proto.Unmarshal(data, pbsensors); err != nil {
		return nil, err
	}
	sensors := make([]*types.Sensor, len(pbsensors.Sensors))
	for i := 0; i < len(pbsensors.Sensors); i++ {
		sensor := new(types.Sensor)
		sensor.FromSnapshot(pbsensors.Sensors[i])
		sensors[i] = sensor
	}
	return sensors, nil
}

// EncodeRoutes encodes a slice of navigator routes into binary data.
func EncodeRoutes(routes []*navigator.Route) ([]byte, error) {
	if len(routes) == 0 {
		return nil, fmt.Errorf("gpsgen: no routes to encode")
	}
	pbroutes := &pb.Snapshot_Navigator_Routes{
		Routes: make([]*pb.Snapshot_Navigator_Route, len(routes)),
	}
	for i := 0; i < len(routes); i++ {
		pbroutes.Routes[i] = routes[i].Snapshot()
	}
	return proto.Marshal(pbroutes)
}

// DecodeRoutes decodes binary data into a slice of navigator routes.
func DecodeRoutes(data []byte) ([]*navigator.Route, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gpsgen: no data to decode routes")
	}
	pbroutes := new(pb.Snapshot_Navigator_Routes)
	if err := proto.Unmarshal(data, pbroutes); err != nil {
		return nil, err
	}
	routes := make([]*navigator.Route, len(pbroutes.Routes))
	for i := 0; i < len(pbroutes.Routes); i++ {
		route := new(navigator.Route)
		route.RouteFromSnapshot(pbroutes.Routes[i])
		routes[i] = route
	}
	return routes, nil
}

// EncodeGeoJSONRoutes encodes a slice of navigator routes into GeoJSON format.
func EncodeGeoJSONRoutes(routes []*navigator.Route) ([]byte, error) {
	return geojson.Encode(routes)
}

// DecodeGeoJSONRoutes decodes GeoJSON data into a slice of navigator routes.
func DecodeGeoJSONRoutes(data []byte) ([]*navigator.Route, error) {
	return geojson.Decode(data)
}

// EncodeGPXRoutes encodes a slice of navigator routes into GPX format.
func EncodeGPXRoutes(routes []*navigator.Route) ([]byte, error) {
	return gpx.Encode(routes)
}

// DecodeGPXRoutes decodes GPX data into a slice of navigator routes.
func DecodeGPXRoutes(data []byte) ([]*navigator.Route, error) {
	return gpx.Decode(data)
}
