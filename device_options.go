package gpsgen

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/random"
	"github.com/mmadfox/go-gpsgen/types"
)

// DeviceOptions defines the options for creating a new device.
type DeviceOptions struct {
	ID     string // ID of the device.
	Model  string // Model of the device.
	Color  string // Color of the device.
	UserID string // User ID associated with the device.
	Descr  string // Description of the device.

	Navigator struct {
		SkipOffline bool     // Skip offline mode.
		Offline     struct { // Offline mode settings.
			Min int // Minimum duration for offline mode.
			Max int // Maximum duration for offline mode.
		}
		Elevation struct { // Elevation settings.
			Min       float64          // Minimum elevation.
			Max       float64          // Maximum elevation.
			Amplitude int              // Amplitude for elevation changes.
			Mode      types.SensorMode // Sensor mode for elevation changes.
		}
	}

	Battery struct { // Battery settings.
		Min        float64       // Minimum battery level.
		Max        float64       // Maximum battery level.
		ChargeTime time.Duration // Charging time for the battery.
	}

	Speed struct { // Speed settings.
		Min       float64 // Minimum speed.
		Max       float64 // Maximum speed.
		Amplitude int     // Amplitude for speed changes.
	}
}

func (o *DeviceOptions) applyFor(s *pb.Device) error {
	if len(o.ID) == 0 {
		o.ID = uuid.NewString()
	}
	s.Id = o.ID
	s.Description = o.Descr
	if len(o.Model) > 0 {
		model, err := types.NewModel(o.Model)
		if err != nil {
			return err
		}
		s.Model = model.String()
	} else {
		s.Model = "Device-" + random.String(8)
	}
	s.UserId = o.UserID
	if len(o.Color) > 0 {
		c, err := colorful.Hex(o.Color)
		if err != nil {
			return err
		}
		s.Color = c.Hex()
	} else {
		s.Color = colorful.FastHappyColor().Hex()
	}
	return nil
}

func (o *DeviceOptions) navOpts() []navigator.Option {
	opts := make([]navigator.Option, 0)
	opts = append(opts,
		navigator.WithElevation(
			o.Navigator.Elevation.Min,
			o.Navigator.Elevation.Max,
			o.Navigator.Elevation.Amplitude,
			o.Navigator.Elevation.Mode,
		))
	opts = append(opts,
		navigator.WithOffline(
			o.Navigator.Offline.Min,
			o.Navigator.Offline.Max,
		))
	if o.Navigator.SkipOffline {
		opts = append(opts, navigator.SkipOfflineMode())
	}
	return opts
}

// NewDeviceOptions creates a new set of default DeviceOptions.
func NewDeviceOptions() *DeviceOptions {
	opts := new(DeviceOptions)

	opts.Navigator.Elevation.Min = 10
	opts.Navigator.Elevation.Max = 300
	opts.Navigator.Elevation.Amplitude = 8

	opts.Navigator.Offline.Min = 2
	opts.Navigator.Offline.Max = 60

	opts.Battery.Min = 0
	opts.Battery.Max = 100
	opts.Battery.ChargeTime = 7 * time.Hour

	opts.Speed.Min = 1
	opts.Speed.Max = 2
	opts.Speed.Amplitude = 16

	return opts
}

// DefaultTrackerOptions returns DeviceOptions suitable for a default tracker.
func DefaultTrackerOptions() *DeviceOptions {
	opts := NewDeviceOptions()
	opts.Descr = "Tracker"
	opts.Navigator.Elevation.Min = 0
	opts.Navigator.Elevation.Max = 350
	opts.Battery.ChargeTime = 8 * time.Hour
	opts.Speed.Min = 1
	opts.Speed.Max = 5
	opts.Speed.Amplitude = 16
	return opts
}

// KidsTrackerOptions returns DeviceOptions suitable for a kids tracker.
func KidsTrackerOptions() *DeviceOptions {
	opts := NewDeviceOptions()
	opts.Descr = "Kids tracker"
	opts.Navigator.Elevation.Min = 0
	opts.Navigator.Elevation.Max = 50
	opts.Battery.ChargeTime = 6 * time.Hour
	opts.Speed.Min = 1
	opts.Speed.Max = 2
	opts.Speed.Amplitude = 48
	return opts
}

// DogTrackerOptions returns DeviceOptions suitable for a dog tracker.
func DogTrackerOptions() *DeviceOptions {
	opts := NewDeviceOptions()
	opts.Descr = "Dog tracker"
	opts.Navigator.Elevation.Min = 1
	opts.Navigator.Elevation.Max = 10
	opts.Battery.ChargeTime = 2 * time.Hour
	opts.Speed.Min = 1
	opts.Speed.Max = 3
	opts.Speed.Amplitude = 64
	return opts
}

// BicycleTrackerOptions returns DeviceOptions suitable for a bicycle tracker.
func BicycleTrackerOptions() *DeviceOptions {
	opts := NewDeviceOptions()
	opts.Descr = "Bicycle tracker"
	opts.Navigator.Elevation.Min = 1
	opts.Navigator.Elevation.Max = 125
	opts.Battery.ChargeTime = 12 * time.Hour
	opts.Speed.Min = 2
	opts.Speed.Max = 8
	opts.Speed.Amplitude = 8
	return opts
}

// DroneTrackerOptions returns DeviceOptions suitable for a drone tracker.
func DroneTrackerOptions() *DeviceOptions {
	opts := NewDeviceOptions()
	opts.Descr = "Drone tracker"
	opts.Navigator.Elevation.Min = 1
	opts.Navigator.Elevation.Max = 2500
	opts.Navigator.SkipOffline = true
	opts.Battery.ChargeTime = 3 * time.Hour
	opts.Speed.Min = 5
	opts.Speed.Max = 150
	opts.Speed.Amplitude = 4
	return opts
}
