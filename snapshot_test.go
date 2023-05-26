package gpsgen

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/route"
	"github.com/stretchr/testify/require"
)

func TestDeviceSnapshot(t *testing.T) {
	routes, err := route.RoutesForChina()
	require.NoError(t, err)
	require.NotNil(t, routes)

	props := Properties{"foo": "bar"}
	sensors := []Sensor{{Name: "sensor", Min: 1, Max: 2, Amplitude: Amplitude16}}
	myDrone, err := DroneWithSensors("myDrone", routes, props, sensors...)
	require.NoError(t, err)
	require.NotNil(t, myDrone)

	snapshot, err := TakeDeviceSnapshot(myDrone)
	require.NoError(t, err)
	require.NotZero(t, snapshot)

	myDrone2, err := DeviceFromSnapshot(snapshot)
	require.NoError(t, err)
	require.NotNil(t, myDrone2)

	for i := 0; i < 3; i++ {
		NextTick(myDrone, myDrone2)
	}

	want := myDrone.State()
	got := myDrone2.State()
	require.Equal(t, want.BatteryCharge, got.BatteryCharge)
	require.Equal(t, want.Descr, got.Descr)
	require.Equal(t, want.Id, got.Id)
	require.Equal(t, want.Location, got.Location)
	require.Equal(t, want.Model, got.Model)
	require.Equal(t, want.Online, got.Online)
	require.Equal(t, want.Props, got.Props)
	require.Equal(t, want.Sensors, got.Sensors)
	require.Equal(t, want.Speed, got.Speed)
	require.Equal(t, want.Tick, got.Tick)
}

func TestDeviceSnapshotError(t *testing.T) {
	snapshot, err := TakeDeviceSnapshot(nil)
	require.Nil(t, snapshot)
	require.Error(t, err)

	dev, err := DeviceFromSnapshot(nil)
	require.Nil(t, dev)
	require.Error(t, err)

	dev, err = DeviceFromSnapshot([]byte("hoho"))
	require.Nil(t, dev)
	require.Error(t, err)
}
