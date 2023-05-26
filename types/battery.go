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

// Battery represents a Battery object with various methods and fields.
type Battery struct {
	min        float64
	max        float64
	val        float64
	chargeTime time.Duration
}

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
	if chargeTime <= 0 {
		chargeTime = defaultChargeTime
	}
	return &Battery{
		min:        min,
		max:        max,
		chargeTime: chargeTime,
	}, nil
}

func (t *Battery) ToProto() *proto.BatteryState {
	return &proto.BatteryState{
		Min:        t.min,
		Max:        t.max,
		ChargeTime: int64(t.chargeTime),
	}
}

func (t *Battery) FromProto(battery *proto.BatteryState) {
	t.min = battery.Min
	t.max = battery.Max
	t.chargeTime = time.Duration(battery.ChargeTime)
}

// Min returns the minimum value of the battery.
func (t *Battery) Min() float64 {
	return t.min
}

// Max returns the maximum value of the battery.
func (t *Battery) Max() float64 {
	return t.max
}

func (t *Battery) ChargeTime() time.Duration {
	return t.chargeTime
}

func (t *Battery) IsLow() bool {
	return t.Value() == t.min
}

// Value returns the current value of the battery.
func (t *Battery) Value() float64 {
	val := 100 - (t.val / t.chargeTime.Seconds() * 100)
	if val < t.min {
		val = t.min
	}
	return val
}

func (t *Battery) Reset() {
	t.val = 0
}

// String returns a formatted string representation of the battery object,
// displaying the battery level as a percentage.
func (t *Battery) String() string {
	return fmt.Sprintf("batteryCharge: %.2f%%", t.Value())
}

func (t *Battery) Next(seconds float64) {
	t.val += seconds
}
