package navigator

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen/types"
)

const (
	MinOffline   = 0
	MaxOffline   = 3600
	MinElevation = 0.0
	MaxElevation = 50000.0
)

// Option is a function type that modifies navigatorOptions.
type Option func(*navigatorOptions)

// SkipOfflineMode skips offline mode in the navigator.
func SkipOfflineMode() Option {
	return func(opt *navigatorOptions) {
		opt.skipOffline = true
	}
}

// WithElevation sets the elevation options for the navigator.
func WithElevation(min, max float64, amplitude int, mode types.SensorMode) Option {
	return func(opt *navigatorOptions) {
		opt.elevationMin, opt.elevationMax = min, max
		opt.elevationAmplitude = amplitude
		opt.elevationMode = mode
	}
}

// WithOffline sets the offline time options for the navigator.
func WithOffline(min, max int) Option {
	return func(opt *navigatorOptions) {
		opt.minOffline, opt.maxOffline = min, max
	}
}

type navigatorOptions struct {
	elevationMin, elevationMax                 float64
	elevationMode                              types.SensorMode
	minOffline, maxOffline, elevationAmplitude int
	skipOffline                                bool
}

func (o *navigatorOptions) validate() error {
	if err := o.validateOffline(); err != nil {
		return err
	}
	return o.validateElevation()
}

func (o *navigatorOptions) validateElevation() error {
	if o.elevationMin < MinElevation {
		return fmt.Errorf("invalid minimum navigator elevation value got %.2f, expected > %.2f meters",
			o.elevationMin, MinElevation)
	}
	if o.elevationMax > MaxElevation {
		return fmt.Errorf("invalid maximum navigator elevation value got %.2f, expected < %.2f meters",
			o.elevationMax, MaxElevation)
	}
	if o.elevationMin > o.elevationMax {
		o.elevationMin = o.elevationMax - 1
	}
	if o.elevationMax < o.elevationMin {
		o.elevationMax = o.elevationMin + 1
	}
	if o.elevationMin == 0 && o.elevationMax == 0 {
		o.elevationMin = MinElevation
		o.elevationMax = o.elevationMin + 100
	}
	return nil
}

func (o *navigatorOptions) validateOffline() error {
	if o.minOffline < MinOffline {
		return fmt.Errorf("invalid minimum navigator offline value got %d, expected > %d seconds",
			o.minOffline, MinOffline)
	}
	if o.maxOffline > MaxOffline {
		return fmt.Errorf("invalid maximum navigator offline value got %d, expected < %d seconds",
			o.maxOffline, MaxOffline)
	}
	if o.minOffline > o.maxOffline {
		o.minOffline = o.maxOffline - 1
	}
	if o.maxOffline < o.minOffline {
		o.maxOffline = o.minOffline + 1
	}
	if o.minOffline == 0 && o.maxOffline == 0 {
		o.minOffline = MinOffline
		o.maxOffline = MinOffline + 10
	}
	return nil
}

func defaultOptions() *navigatorOptions {
	return &navigatorOptions{
		elevationMin:       1,
		elevationMax:       500,
		elevationAmplitude: 4,
		minOffline:         1,
		maxOffline:         10,
	}
}
