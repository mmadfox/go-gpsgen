package route

import (
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"google.golang.org/protobuf/proto"
)

// Encode converts a slice of navigator.Route objects to a
// binary representation using Protocol Buffers (proto.Rotues).
func Encode(routes []*navigator.Route) ([]byte, error) {
	protoRoutes := &pb.Routes{
		Routes: make([]*pb.Route, len(routes)),
	}
	for i := 0; i < len(routes); i++ {
		protoRoutes.Routes[i] = routes[i].ToProto()
	}
	return proto.Marshal(protoRoutes)
}

// Decode converts a byte array representing a
// binary encoded protobuf message back into a slice of navigator.Route objects.
func Decode(data []byte) ([]*navigator.Route, error) {
	protoRoutes := new(pb.Routes)
	if err := proto.Unmarshal(data, protoRoutes); err != nil {
		return nil, err
	}
	routes := make([]*navigator.Route, len(protoRoutes.Routes))
	for i := 0; i < len(protoRoutes.Routes); i++ {
		route := new(navigator.Route)
		route.FromProto(protoRoutes.Routes[i])
		routes[i] = route
	}
	return routes, nil
}
