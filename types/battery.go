package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmadfox/go-gpsgen/proto"
)

const (
	minBatteryVal     = 0
	maxBatteryVal     = 100
	defaultChargeTime = time.Hour * 9
)

var (
	// ErrMinBattery indicates that the value is less than 0%.
	ErrMinBattery = errors.New("types/battery: value is less than 0%")
	// ErrMaxBattery indicates that the value greater than 100%.
	ErrMaxBattery = errors.New("types/battery: value greater than 100%")
)

// Battery represents a simulated battery object with charge level management and behavior.
type Battery struct {
	min        float64
	max        float64
	val        float64
	chargeTime time.Duration
}

// NewBattery creates a new Battery instance with the specified characteristics.
// The min parameter is the minimum battery charge level (percentage),
// the max parameter is the maximum battery charge level (percentage),
// and the chargeTime parameter represents the time duration it takes to fully charge the battery.
// If the min value is less than 0 or the max value is greater than 100, an error is returned.
// If the min value is greater than the max value, the min value is adjusted to match the max value.
// If the chargeTime is less than or equal to zero, it is set to a default value of 9 hours.
func NewBattery(min, max float64, chargeTime time.Duration) (*Battery, error) {
	if min < minBatteryVal {
		return nil, ErrMinBattery
	}
	if max > maxBatteryVal {
		return nil, ErrMaxBattery
	}
	if min > max {
		min = max
	}
	if chargeTime < 0 {
		chargeTime = defaultChargeTime
	}
	return &Battery{
		min:        min,
		max:        max,
		chargeTime: chargeTime,
	}, nil
}

// Min returns the minimum value of the battery.
func (t *Battery) Min() float64 {
	return t.min
}

// Max returns the maximum value of the battery.
func (t *Battery) Max() float64 {
	return t.max
}

// ChargeTime returns the time duration it takes to fully charge the battery.
func (t *Battery) ChargeTime() time.Duration {
	return t.chargeTime
}

// IsLow returns true if the battery charge level is at its minimum value, indicating a low battery condition.
func (t *Battery) IsLow() bool {
	val := t.Value()
	return val == t.min
}

// Value returns the current value of the battery.
func (t *Battery) Value() float64 {
	val := 100 - (t.val / t.chargeTime.Seconds() * 100)
	if val < t.min {
		val = t.min
	}
	return val
}

// Reset resets the battery charge level to zero.
func (t *Battery) Reset() {
	t.val = 0
}

// String returns a formatted string representation of the battery object,
// displaying the battery level as a percentage.
func (t *Battery) String() string {
	return fmt.Sprintf("batteryCharge: %.2f%%", t.Value())
}

// Next updates the battery charge level based on the number of seconds that have passed.
func (t *Battery) Next(seconds float64) {
	t.val += seconds
}

// Snapshot returns a protobuf snapshot of the Battery instance.
func (t *Battery) Snapshot() *proto.Snapshot_BatteryType {
	return &proto.Snapshot_BatteryType{
		Min:        t.min,
		Max:        t.max,
		Val:        t.val,
		ChargeTime: int64(t.chargeTime),
	}
}

// FromSnapshot restores the Battery instance from a protobuf snapshot.
func (t *Battery) FromSnapshot(snap *proto.Snapshot_BatteryType) {
	t.min = snap.Min
	t.max = snap.Max
	t.val = snap.Val
	t.chargeTime = time.Duration(snap.ChargeTime)
}
