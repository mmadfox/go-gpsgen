package gpsgen

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestPacketFromBytes(t *testing.T) {
	myTracker := NewDogTracker()

	route := navigator.RouteFromTracks(track1km2segment)
	myTracker.AddRoute(route)

	state := myTracker.State()
	expectedPacket := &pb.Packet{
		Devices: []*pb.Device{state},
	}

	data, err := proto.Marshal(expectedPacket)
	require.NoError(t, err)

	actualPacket, err := PacketFromBytes(data)
	require.NotNil(t, actualPacket)
	require.NoError(t, err)
	require.Equal(t, len(expectedPacket.Devices), len(actualPacket.Devices))

	for i := 0; i < len(expectedPacket.Devices); i++ {
		s1 := expectedPacket.Devices[i]
		s2 := actualPacket.Devices[i]
		require.Equal(t, s1.Navigator, s2.Navigator)
	}
}

func TestPacketFromBytesWithNilData(t *testing.T) {
	packet, err := PacketFromBytes(nil)
	require.Error(t, err)
	require.Nil(t, packet)
}

func TestPacketFromBytesWithInvalidData(t *testing.T) {
	packet, err := PacketFromBytes([]byte("somedata"))
	require.Error(t, err)
	require.Nil(t, packet)
}
