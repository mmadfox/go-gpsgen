package types

import "errors"

const (
	minAmplitude = 4
	maxAmplitude = 512
)

var (
	// ErrMinAmplitude indicates that the amplitude value is less than 4.
	ErrMinAmplitude = errors.New("types/amplitude: value is less than 4")
	// ErrMaxAmplitude indicates that the amplitude value of greater 512.
	ErrMaxAmplitude = errors.New("types/amplitude: value of greater 512")
)

func validateAmplitude(val int) error {
	if val < minAmplitude {
		return ErrMinAmplitude
	}
	if val > maxAmplitude {
		return ErrMaxAmplitude
	}
	return nil
}

func validateSensorMode(mode SensorMode) bool {
	return mode > 0
}
