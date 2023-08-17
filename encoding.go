package gpsgen

import (
	"fmt"

	pb "github.com/mmadfox/go-gpsgen/proto"
	"google.golang.org/protobuf/proto"
)

// PacketFromBytes decodes a byte slice into a protobuf Packet.
// Returns the decoded Packet and an error if decoding fails.
func PacketFromBytes(data []byte) (*pb.Packet, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gpsgen: invalid packet size")
	}
	pck := new(pb.Packet)
	if err := proto.Unmarshal(data, pck); err != nil {
		return nil, err
	}
	return pck, nil
}
