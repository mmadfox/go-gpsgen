package gpsgen

import (
	"fmt"
	"math"
	"sync"

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
	id            uuid.UUID
	userID        string
	props         Properties
	descr         string
	model         types.Model
	speed         *types.Speed
	battery       *types.Battery
	sensors       []*types.Sensor
	navigator     *navigator.Navigator
	OnStateChange func(*State, []byte)
	stateCh       chan *State
	readyCh       chan struct{}
	loop          float64
	avgTicks      float64
	pool          *sync.Pool
}

// The State struct represents the state of a device.
// It contains fields that provide information about the device's state.
type State struct {
	ID       uuid.UUID             `json:"id"`
	Tick     float64               `json:"tick"`
	UserID   string                `json:"userID"`
	Model    string                `json:"model"`
	Speed    float64               `json:"speed"`
	Battery  float64               `json:"battery"`
	Sensors  map[string][2]float64 `json:"sensors"`
	Location navigator.Location    `json:"location"`
	Props    Properties            `json:"props"`
	Descr    string                `json:"descr"`
	Online   bool                  `json:"online"`
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
		stateCh:   make(chan *State, 1),
		readyCh:   make(chan struct{}, 1),
		pool: &sync.Pool{
			New: func() any {
				return &State{
					Sensors: make(map[string][2]float64, 0),
				}
			},
		},
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
func (d *Device) State() *State {
	state := new(State)
	d.fillState(state)
	return state
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
	d.stateCh = make(chan *State, 1)
	d.readyCh = make(chan struct{}, 1)
	d.pool = &sync.Pool{
		New: func() any {
			return &State{
				Sensors: make(map[string][2]float64, 0),
			}
		},
	}
	return nil
}

func (d *Device) handleChange() {
	snapshot := &pb.Device{}
	for state := range d.stateCh {
		if d.OnStateChange == nil {
			continue
		}
		d.fillSnapshot(state, snapshot)
		data, _ := proto.Marshal(snapshot)
		d.OnStateChange(state, data)
		d.pool.Put(state)
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
		state, ok := d.pool.Get().(*State)
		if !ok {
			state = d.pool.New().(*State)
		}
		d.fillState(state)
		select {
		case d.stateCh <- state:
		default:
		}
	}

	return next
}

func (d *Device) fillState(state *State) {
	state.Model = d.model.String()
	state.Battery = d.battery.Value()
	state.Speed = d.speed.Value()
	state.Tick = d.loop
	state.Online = d.navigator.IsOnline()
	state.Location = d.navigator.Location()
	if state.Sensors == nil {
		state.Sensors = make(map[string][2]float64, len(d.sensors))
	}
	if d.props != nil {
		if state.Props == nil {
			state.Props = make(Properties, len(d.props))
		}
		for k, v := range d.props {
			state.Props[k] = v
		}
	}
	for i := 0; i < len(d.sensors); i++ {
		sensor := d.sensors[i]
		state.Sensors[sensor.Name()] = [2]float64{
			sensor.ValueX(),
			sensor.ValueY(),
		}
	}
	if len(state.UserID) == 0 {
		state.UserID = d.userID
	}
	if len(state.Descr) == 0 {
		state.Descr = d.descr
	}
}

func (d *Device) fillSnapshot(state *State, snap *pb.Device) {
	snap.Battery = state.Battery
	snap.Speed = state.Speed
	snap.Tick = int64(state.Tick)
	snap.Online = state.Online
	snap.Latitude = state.Location.Lat
	snap.Longitude = state.Location.Lon
	snap.Elevation = state.Location.Alt
	snap.CurrentDistance = state.Location.CurrentDistance
	snap.TotalDistance = state.Location.TotalDistance

	if snap.Lat == nil {
		snap.Lat = new(pb.DMS)
	}
	if snap.Lon == nil {
		snap.Lon = new(pb.DMS)
	}
	if snap.Utm == nil {
		snap.Utm = new(pb.UTM)
	}

	snap.Lat.Degrees = int64(state.Location.LatDMS.Degrees)
	snap.Lat.Direction = state.Location.LatDMS.Direction
	snap.Lat.Minutes = int64(state.Location.LatDMS.Minutes)
	snap.Lat.Seconds = state.Location.LatDMS.Seconds

	snap.Lon.Degrees = int64(state.Location.LonDMS.Degrees)
	snap.Lon.Direction = state.Location.LonDMS.Direction
	snap.Lon.Minutes = int64(state.Location.LonDMS.Minutes)
	snap.Lon.Seconds = state.Location.LonDMS.Seconds

	snap.Utm.CentralMeridian = state.Location.CurrentDistance
	snap.Utm.Easting = state.Location.UTM.Easting
	snap.Utm.Hemisphere = state.Location.UTM.Hemisphere
	snap.Utm.LatZone = state.Location.UTM.LatZone
	snap.Utm.LongZone = int64(state.Location.UTM.LongZone)
	snap.Utm.Northing = state.Location.UTM.Northing
	snap.Utm.Srid = int64(state.Location.UTM.SRID)

	if snap.Sensors == nil && len(state.Sensors) > 0 {
		snap.Sensors = make([]*pb.Sensor, 0, len(state.Sensors))
		for name := range state.Sensors {
			snap.Sensors = append(snap.Sensors, &pb.Sensor{
				Name: name,
			})
		}
	}
	for i := 0; i < len(snap.Sensors); i++ {
		sensor := snap.Sensors[i]
		stateSensor, ok := state.Sensors[sensor.Name]
		if !ok {
			continue
		}
		sensor.ValX = stateSensor[0]
		sensor.ValY = stateSensor[1]
	}
	if len(snap.Model) == 0 {
		snap.Model = state.Model
	}
	if len(snap.Props) == 0 && len(state.Props) > 0 {
		snap.Props = make(map[string]string, len(state.Props))
		for k, v := range state.Props {
			snap.Props[k] = v
		}
	}
	if len(snap.Descr) == 0 && len(state.Descr) > 0 {
		snap.Descr = state.Descr
	}
}
