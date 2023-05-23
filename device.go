package gpsgen

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/types"
	"google.golang.org/protobuf/proto"
)

// The Device struct represents a device and contains information
// about the device's properties, description, model, speed,
// battery, sensors, navigator, and various other fields.
type Device struct {
	id        uuid.UUID
	userID    string
	props     Properties
	descr     string
	model     types.Model
	speed     *types.Speed
	battery   *types.Battery
	sensors   []*types.Sensor
	navigator *navigator.Navigator
	stateCh   chan struct{}
	readyCh   chan struct{}
	loop      float64
	avgTicks  float64
	state     *pb.Device

	OnStateChange      func(*pb.Device)
	OnStateChangeBytes func([]byte)
}

// Properties describes custom device characteristics.
type Properties map[string]string

// NewDevice creates a new instance of the Device struct with the given settings.
// It applies the provided settings to the device and returns a pointer to the created Device and an error, if any.
func NewDevice(settings ...DeviceSetting) (*Device, error) {
	opts := defaultSettings()
	for i := 0; i < len(settings); i++ {
		settings[i](opts)
	}

	nav, err := navigator.New(
		navigator.WithElevation(
			opts.elevation.min,
			opts.elevation.max,
			float64(opts.elevation.amplitude),
		),
		navigator.WithOffline(opts.offline.min, opts.offline.max),
	)
	if err != nil {
		return nil, err
	}

	deviceID := uuid.New()

	device := &Device{
		id:        deviceID,
		userID:    opts.userID,
		props:     opts.props,
		descr:     opts.descr,
		navigator: nav,
		stateCh:   make(chan struct{}, 1),
		readyCh:   make(chan struct{}, 1),
		state: &pb.Device{
			Model:  opts.model,
			Descr:  opts.descr,
			UserId: opts.userID,
			Location: &pb.Location{
				LatDms: new(pb.DMS),
				LonDms: new(pb.DMS),
				Utm:    new(pb.UTM),
			},
		},
	}

	if len(opts.sensors) > 0 {
		device.state.Sensors = make([]*pb.Sensor, len(opts.sensors))
		for i := 0; i < len(opts.sensors); i++ {
			device.state.Sensors[i] = &pb.Sensor{
				Name: opts.sensors[i].Name,
			}
		}
	}

	if len(opts.props) > 0 {
		device.state.Props = make(map[string]string, len(opts.props))
		for k, v := range opts.props {
			device.state.Props[k] = v
		}
	}

	if err := opts.applyFor(device); err != nil {
		return nil, err
	}

	avgSpeed := (device.speed.Min() + device.speed.Max()) / 2
	device.avgTicks = device.navigator.TotalDistance() / avgSpeed

	return device, nil
}

// ID returns the unique identifier of the device.
func (d *Device) ID() uuid.UUID {
	return d.id
}

// State returns a State object filled with the current state of the device.
func (d *Device) State() *pb.Device {
	d.fillState()
	return d.state
}

// MarshalBinary serializes the Device struct into a binary representation using Protocol Buffers.
// The method marshals the proto.DeviceSnapshot object into a byte slice using Protocol Buffers' proto.Marshal function.
func (d *Device) MarshalBinary() ([]byte, error) {
	protoDev := &pb.DeviceSnapshot{
		Id:          d.id[:],
		Model:       d.model.String(),
		Speed:       d.speed.ToProto(),
		Battery:     d.battery.ToProto(),
		Sensors:     make([]*pb.SensorState, len(d.sensors)),
		Navigator:   d.navigator.ToProto(),
		Loop:        d.loop,
		AvgTick:     d.avgTicks,
		UserId:      d.userID,
		Description: d.descr,
		Properties:  make(Properties, len(d.props)),
	}
	for k, v := range d.props {
		protoDev.Properties[k] = v
	}
	for i := 0; i < len(d.sensors); i++ {
		protoDev.Sensors[i] = d.sensors[i].ToProto()
	}
	return proto.Marshal(protoDev)
}

// UnmarshalBinary deserializes a binary representation of a Device from a byte slice.
// The method unmarshals the binary data into a proto.DeviceSnapshot object and assigns
// the values to the corresponding fields of the Device struct.
func (d *Device) UnmarshalBinary(data []byte) error {
	protoDev := new(pb.DeviceSnapshot)
	if err := proto.Unmarshal(data, protoDev); err != nil {
		return err
	}
	if protoDev.Navigator == nil {
		return fmt.Errorf("invalid device snapshot data")
	}
	d.id = uuid.UUID(protoDev.Id)
	d.model, _ = types.NewModel(protoDev.Model)
	d.speed = new(types.Speed)
	d.speed.FromProto(protoDev.Speed)
	d.battery = new(types.Battery)
	d.battery.FromProto(protoDev.Battery)
	d.sensors = make([]*types.Sensor, len(protoDev.Sensors))
	for i := 0; i < len(protoDev.Sensors); i++ {
		d.sensors[i] = new(types.Sensor)
		d.sensors[i].FromProto(protoDev.Sensors[i])
	}
	d.navigator = new(navigator.Navigator)
	d.navigator.FromProto(protoDev.Navigator)
	d.loop = protoDev.Loop
	d.avgTicks = protoDev.AvgTick
	d.userID = protoDev.UserId
	d.props = make(Properties, len(protoDev.Properties))
	for k, v := range protoDev.Properties {
		d.props[k] = v
	}
	d.descr = protoDev.Description
	d.stateCh = make(chan struct{}, 1)
	d.readyCh = make(chan struct{}, 1)
	d.state = &pb.Device{
		Model:  protoDev.Model,
		Descr:  protoDev.Description,
		UserId: protoDev.UserId,
		Location: &pb.Location{
			LatDms: new(pb.DMS),
			LonDms: new(pb.DMS),
			Utm:    new(pb.UTM),
		},
	}

	if len(protoDev.Sensors) > 0 {
		d.state.Sensors = make([]*pb.Sensor, len(protoDev.Sensors))
		for i := 0; i < len(protoDev.Sensors); i++ {
			d.state.Sensors[i] = &pb.Sensor{
				Name: protoDev.Sensors[i].Name,
			}
		}
	}

	if len(protoDev.Properties) > 0 {
		d.state.Props = make(map[string]string, len(protoDev.Properties))
		for k, v := range protoDev.Properties {
			d.state.Props[k] = v
		}
	}

	return nil
}

func (d *Device) handleChange() {
	for {
		<-d.stateCh
		if d.OnStateChange != nil {
			d.OnStateChange(d.state)
		}
		if d.OnStateChangeBytes != nil {
			data, _ := proto.Marshal(d.state)
			d.OnStateChangeBytes(data)
		}
		d.readyCh <- struct{}{}
	}
}

func (d *Device) nextTick(tick float64) bool {
	if !d.navigator.IsOnline() {
		d.navigator.NextOffline()
		return false
	}

	var isReady bool
	select {
	case <-d.readyCh:
		isReady = true
	default:
		if d.loop == 0 {
			isReady = true
		}
	}

	d.loop += tick
	t := math.Min(d.loop/d.avgTicks, 1.0)
	d.speed.Next(t)

	d.battery.Next(t)
	d.navigator.NextSensors(t)

	if len(d.sensors) > 0 {
		for i := 0; i < len(d.sensors); i++ {
			d.sensors[i].Next(t)
		}
	}

	var finish bool
	if d.navigator.IsFinish() {
		finish = true
		d.navigator.ToOffline()
	}

	next := d.navigator.Next(tick, d.speed.Value())
	if !next && finish {
		d.loop = 0
	}

	if isReady && d.navigator.CurrentDistance() > 0 {
		d.fillState()
		select {
		case d.stateCh <- struct{}{}:
		default:
		}
	}

	return next
}

func (d *Device) fillState() {
	d.state.Battery = d.battery.Value()
	d.state.Speed = d.speed.Value()
	d.state.Tick = int64(d.loop)
	d.state.Online = d.navigator.IsOnline()
	d.navigator.UpdateLocation(d.state.Location)

	if len(d.sensors) > 0 {
		for i := 0; i < len(d.sensors); i++ {
			d.state.Sensors[i].ValX = d.sensors[i].ValueX()
			d.state.Sensors[i].ValY = d.sensors[i].ValueY()
		}
	}
}
