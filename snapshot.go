package gpsgen

import "fmt"

// ToDeviceSnapshot converts a Device object into a byte array that represents a
// snapshot of the device's state.
func ToDeviceSnapshot(d *Device) ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("device is nil pointer")
	}
	return d.MarshalBinary()
}

// FromDeviceSnapshot restores a Device object from the passed byte array representing a
// snapshot of the device's state.
func FromDeviceSnapshot(state []byte) (*Device, error) {
	dev := new(Device)
	if err := dev.UnmarshalBinary(state); err != nil {
		return nil, err
	}
	return dev, nil
}
