package gpsgen

import (
	"fmt"
	"time"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/mmadfox/go-gpsgen/types"
)

// ConfigurationError type represents an error specific to configuration issues.
type ConfigurationError string

// FormatConfigurationError formats an error message with the given
// format string and arguments and returns a ConfigurationError instance.
func FormatConfigurationError(format string, a ...any) ConfigurationError {
	return ConfigurationError(fmt.Sprintf(format, a...))
}

// Error returns the error message as a string.
func (err ConfigurationError) Error() string {
	return "gpsgen: invalid configuration (" + string(err) + ")"
}

// Config struct represents a device configuration.
type Config struct {
	Model       string            `json:"model"`
	UserID      string            `json:"userId,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	Description string            `json:"description,omitempty"`
	Speed       struct {
		Max       float64 `json:"max"`
		Min       float64 `json:"min"`
		Amplitude int     `json:"amplitude"`
	} `json:"speed"`
	Battery struct {
		Max        float64       `json:"max"`
		Min        float64       `json:"min"`
		ChargeTime time.Duration `json:"chargeTime"`
	} `json:"battery"`
	Elevation struct {
		Max       float64 `json:"max"`
		Min       float64 `json:"min"`
		Amplitude int     `json:"amplitude"`
	} `json:"elevation"`
	Offline struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"offline"`
	Sensors []Sensor           `json:"sensors,omitempty"`
	Routes  []*navigator.Route `json:"-"`

	isValid bool
}

// NewConfig creates and initializes a new Config instance with default values.
func NewConfig() *Config {
	conf := new(Config)
	conf.Model = types.RandomModel().String()
	conf.Speed.Min = 3
	conf.Speed.Max = 5
	conf.Speed.Amplitude = Amplitude16
	conf.Battery.Min = 1
	conf.Battery.Max = 100
	conf.Battery.ChargeTime = 4 * time.Hour
	conf.Elevation.Min = 5
	conf.Elevation.Max = 130
	conf.Elevation.Amplitude = Amplitude32
	conf.Offline.Min = 1
	conf.Offline.Max = 3
	routes, _ := route.RoutesForRussia()
	conf.Routes = routes
	conf.Description = "drone"
	return conf
}

// NewDevice creates a new Device instance based on the configuration.
// It performs validation on the configuration before creating the device.
func (c *Config) NewDevice() (*Device, error) {
	if c == nil {
		return nil, fmt.Errorf("<nil> pointer")
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewDevice(
		WithModel(c.Model),
		WithUserID(c.UserID),
		WithRoutes(c.Routes),
		WithSpeed(c.Speed.Min, c.Speed.Max, c.Speed.Amplitude),
		WithBattery(c.Battery.Min, c.Battery.Max, c.Battery.ChargeTime),
		WithSensors(c.Sensors...),
		WithElevation(c.Elevation.Min, c.Elevation.Max, c.Elevation.Amplitude),
		WithOffline(c.Offline.Min, c.Offline.Max),
		WithProps(c.Properties),
		WithDescription(c.Description),
	)
}

var constraints = DeviceConstraints()

// Validate validates the configuration against a set of constraints.
// It checks various fields of the configuration for compliance with the defined constraints.
func (c *Config) Validate() error {
	if c.isValid {
		return nil
	}
	funcs := []func() error{
		c.validateModel,
		c.validateProperties,
		c.validateDescription,
		c.validateSpeed,
		c.validateBattery,
		c.validateElevation,
		c.validateSensors,
		c.validateOffline,
		c.validateRoutes,
	}
	for i := 0; i < len(funcs); i++ {
		if err := funcs[i](); err != nil {
			return err
		}
	}
	c.isValid = true
	return nil
}

func (c *Config) validateModel() error {
	if len(c.Model) < constraints.Model.Min {
		return FormatConfigurationError("Model must be > %d", constraints.Model.Min)
	}
	if len(c.Model) > constraints.Model.Max {
		return FormatConfigurationError("Model must be < %d", constraints.Model.Max)
	}
	return nil
}

func (c *Config) validateProperties() error {
	if len(c.Properties) < constraints.Properties.Min {
		return FormatConfigurationError("Properties must be > %d", constraints.Properties.Min)
	}
	if len(c.Properties) > constraints.Properties.Max {
		return FormatConfigurationError("Properties must be < %d", constraints.Properties.Max)
	}
	for k, v := range c.Properties {
		if len(k) > constraints.Properties.MaxKeyLen {
			return FormatConfigurationError("Properties.key{%s} must be < %d", k, constraints.Properties.MaxKeyLen)
		}
		if len(v) > constraints.Properties.MaxValueLen {
			return FormatConfigurationError("Properties.value{%s} must be < %d", v, constraints.Properties.MaxValueLen)
		}
	}
	return nil
}

func (c *Config) validateDescription() error {
	if len(c.Description) == 0 {
		return nil
	}
	if len(c.Description) < constraints.Description.Min {
		return FormatConfigurationError("Description must be > %d chars", constraints.Description.Min)
	}
	if len(c.Description) > constraints.Description.Max {
		return FormatConfigurationError("Description must be < %d chars", constraints.Description.Max)
	}
	return nil
}

func (c *Config) validateSpeed() error {
	if c.Speed.Min < constraints.Speed.Min {
		return FormatConfigurationError("Speed.Min must be > %f", constraints.Speed.Min)
	}
	if c.Speed.Max > constraints.Speed.Max {
		return FormatConfigurationError("Speed.Max must be < %f", constraints.Speed.Max)
	}
	if c.Speed.Amplitude < Amplitude4 {
		return FormatConfigurationError("Speed.Amplitude must be > %d", Amplitude4)
	}
	if c.Speed.Amplitude > Amplitude512 {
		return FormatConfigurationError("Speed.Amplitude must be < %d", Amplitude512)
	}
	return nil
}

func (c *Config) validateBattery() error {
	if c.Battery.Min < constraints.Battery.Min {
		return FormatConfigurationError("Battery.Min must be > %f", constraints.Battery.Min)
	}
	if c.Battery.Max > constraints.Battery.Max {
		return FormatConfigurationError("Battery.Max must be < %f", constraints.Battery.Max)
	}
	return nil
}

func (c *Config) validateElevation() error {
	if c.Elevation.Min < constraints.Elevation.Min {
		return FormatConfigurationError("Elevation.Min must be > %f", constraints.Speed.Min)
	}
	if c.Speed.Max > constraints.Speed.Max {
		return FormatConfigurationError("Elevation.Max must be < %f", constraints.Speed.Max)
	}
	if c.Speed.Amplitude < Amplitude4 {
		return FormatConfigurationError("Elevation.Amplitude must be > %d", Amplitude4)
	}
	if c.Speed.Amplitude > Amplitude512 {
		return FormatConfigurationError("Elevation.Amplitude must be < %d", Amplitude512)
	}
	return nil
}

func (c *Config) validateSensors() error {
	if len(c.Sensors) > constraints.Sensors.Max {
		return FormatConfigurationError("Sensors.Max must be < %d", constraints.Sensors.Max)
	}
	for i := 0; i < len(c.Sensors); i++ {
		if c.Sensors[i].Amplitude < Amplitude4 {
			return FormatConfigurationError("Sensors[%d].Amplitude must be > %d", i, Amplitude4)
		}
		if c.Sensors[i].Amplitude > Amplitude512 {
			return FormatConfigurationError("Sensors[%d].Amplitude must be < %d", i, Amplitude512)
		}
	}
	return nil
}

func (c *Config) validateOffline() error {
	if c.Offline.Max > constraints.Offline.Max {
		return FormatConfigurationError("Offline.Max must be < %d", constraints.Offline.Max)
	}
	return nil
}

func (c *Config) validateRoutes() error {
	if len(c.Routes) < constraints.Routes.Min {
		return FormatConfigurationError("Routes.Min must be > %d", constraints.Routes.Min)
	}
	if len(c.Routes) > constraints.Routes.Max {
		return FormatConfigurationError("Routes.Max must be < %d", constraints.Routes.Max)
	}
	var distance float64
	for i := 0; i < len(c.Routes); i++ {
		route := c.Routes[i]
		if route.NumTracks() < constraints.Routes.MinTracksPerRoute {
			return FormatConfigurationError("Routes[%d].MinTrackPerRoute must be > %d", i, constraints.Routes.MinTracksPerRoute)
		}
		if route.NumTracks() > constraints.Routes.MaxTracksPerRoute {
			return FormatConfigurationError("Routes[%d].MaxTracksPerRoute must be < %d", i, constraints.Routes.MaxTracksPerRoute)
		}
		for j := 0; j < route.NumTracks(); j++ {
			numSegments := route.NumSegments(j)
			if numSegments < constraints.Routes.MinSegmentsPerTrack {
				return FormatConfigurationError("Routes[%d].Tracks[%d].MinSegmentsPerTrack must be > %d", i, j, constraints.Routes.MinSegmentsPerTrack)
			}
			if numSegments > constraints.Routes.MaxSegmentsPerTrack {
				return FormatConfigurationError("Routes[%d].Tracks[%d].MaxSegmentsPerTrack must be < %d", i, j, constraints.Routes.MaxSegmentsPerTrack)
			}
		}
		distance += c.Routes[i].Distance()
	}
	if distance < constraints.Routes.MinDistance {
		return FormatConfigurationError("Routes.MinDistance must be > %f", constraints.Routes.MinDistance)
	}
	if distance > constraints.Routes.MaxDistance {
		return FormatConfigurationError("Routes.MaxDistance must be < %f", constraints.Routes.MaxDistance)
	}
	return nil
}

// DeviceConstraints returns a set of predefined constraints for the device configuration.
// It defines limits and ranges for various properties of the device.
func DeviceConstraints() Constraints {
	c := Constraints{}
	c.Model.Min = types.MinModelLen
	c.Model.Max = types.MaxModelLen
	c.Properties.Min = 0
	c.Properties.Max = 16
	c.Properties.MaxKeyLen = 32
	c.Properties.MaxValueLen = 64
	c.Description.Min = 3
	c.Description.Max = 256
	c.Speed.Min = types.MinSpeedVal
	c.Speed.Max = types.MaxSpeedVal
	c.Speed.AmplitudeMin = Amplitude4
	c.Speed.AmplitudeMax = Amplitude512
	c.Battery.Min = 0
	c.Battery.Max = 100
	c.Elevation.Min = 0
	c.Elevation.Max = 10000
	c.Elevation.AmplitudeMin = Amplitude4
	c.Elevation.AmplitudeMax = Amplitude512
	c.Sensors.Min = 0
	c.Sensors.Max = 15
	c.Routes.Min = 1
	c.Routes.Max = 10
	c.Routes.MinTracksPerRoute = 1
	c.Routes.MaxTracksPerRoute = 128
	c.Routes.MinSegmentsPerTrack = 1
	c.Routes.MaxSegmentsPerTrack = 1000
	c.Routes.MinDistance = 1000
	c.Routes.MaxDistance = 2000000
	c.Offline.Min = 0
	c.Offline.Max = 900
	return c
}

// Constraints struct defines the constraints for different properties of the device configuration.
// It contains fields specifying the minimum and maximum values, lengths, and ranges for different properties.
type Constraints struct {
	Model struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"model"`
	Properties struct {
		Min         int `json:"min"`
		Max         int `json:"max"`
		MaxValueLen int `json:"maxValueLen"`
		MaxKeyLen   int `json:"minKeyLen"`
	} `json:"properties"`
	Description struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"description"`
	Speed struct {
		Max          float64 `json:"max"`
		Min          float64 `json:"min"`
		AmplitudeMin int     `json:"amplitudeMin"`
		AmplitudeMax int     `json:"amplitudeMax"`
	} `json:"speed"`
	Battery struct {
		Max float64 `json:"max"`
		Min float64 `json:"min"`
	} `json:"battery"`
	Elevation struct {
		Max          float64 `json:"max"`
		Min          float64 `json:"min"`
		AmplitudeMin int     `json:"amplitudeMin"`
		AmplitudeMax int     `json:"amplitudeMax"`
	} `json:"elevation"`
	Offline struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"offline"`
	Sensors struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"sensors"`
	Routes struct {
		Min                 int     `json:"min"`
		Max                 int     `json:"max"`
		MinTracksPerRoute   int     `json:"minTracksPerRoute"`
		MaxTracksPerRoute   int     `json:"maxTracksPerRoute"`
		MinSegmentsPerTrack int     `json:"minSegmentsPerTrack"`
		MaxSegmentsPerTrack int     `json:"maxSegmentsPerTrack"`
		MinDistance         float64 `json:"minDistance"`
		MaxDistance         float64 `json:"maxDistance"`
	} `json:"routes"`
}
