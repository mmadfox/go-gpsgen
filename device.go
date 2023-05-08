package gpsgen

import (
	"math"
	"sync"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/types"
	"google.golang.org/protobuf/proto"
)

type Device struct {
	id            uuid.UUID
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
		navigator: nav,
		stateCh:   make(chan *State, 1),
		readyCh:   make(chan struct{}, 1),
		pool: &sync.Pool{
			New: func() any {
				return &State{
					ID:      deviceID,
					UserID:  opts.userID,
					Props:   opts.props,
					Descr:   opts.descr,
					Sensors: make(map[string][2]float64, len(opts.sensors)),
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

type Properties map[string]string

func (d *Device) ID() uuid.UUID {
	return d.id
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
	state.Battery = d.battery.Value()
	state.Speed = d.speed.Value()
	state.Tick = d.loop
	state.Online = d.navigator.IsOnline()
	state.Location = d.navigator.Location()
	for i := 0; i < len(d.sensors); i++ {
		sensor := d.sensors[i]
		state.Sensors[sensor.Name()] = [2]float64{
			sensor.ValueX(),
			sensor.ValueY(),
		}
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
