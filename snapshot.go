package gpsgen

import "fmt"

// TakeDeviceSnapshot converts a Device object into a byte array that represents a
// snapshot of the device's state.
func TakeDeviceSnapshot(d *Device) ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("device is nil pointer")
	}
	return d.MarshalBinary()
}

// DeviceFromSnapshot restores a Device object from the passed byte array representing a
// snapshot of the device's state.
func DeviceFromSnapshot(state []byte) (*Device, error) {
	if len(state) == 0 {
		return nil, fmt.Errorf("invalid device state")
	}
	dev := new(Device)
	if err := dev.UnmarshalBinary(state); err != nil {
		return nil, err
	}
	return dev, nil
}
