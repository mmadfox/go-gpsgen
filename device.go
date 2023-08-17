package gpsgen

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/types"
	"google.golang.org/protobuf/proto"
)

// Status represents the status of a device.
type Status byte

// Running and Stopped are the possible values for the Status type.
const (
	Running Status = iota + 1
	Stopped
)

// Device represents a GPS tracking device with various capabilities.
type Device struct {
	status    Status
	navigator *navigator.Navigator
	speed     *types.Speed
	battery   *types.Battery
	mu        sync.RWMutex
	state     *pb.Device
	sensors   []*types.Sensor
	avgT      float64
	tick      uint32
	ns        [3]int
}

// NewDevice creates a new GPS tracking device with the provided options.
func NewDevice(opts *DeviceOptions) (*Device, error) {
	if opts == nil {
		opts = NewDeviceOptions()
	}

	navigator, err := navigator.New(opts.navOpts()...)
	if err != nil {
		return nil, err
	}

	speed, err := types.NewSpeed(opts.Speed.Min, opts.Speed.Max, opts.Speed.Amplitude)
	if err != nil {
		return nil, err
	}

	battery, err := types.NewBattery(opts.Battery.Min, opts.Battery.Max, opts.Battery.ChargeTime)
	if err != nil {
		return nil, err
	}

	state := initDeviceState()
	if err := opts.applyFor(state); err != nil {
		return nil, err
	}

	dev := &Device{
		status:    Stopped,
		navigator: navigator,
		speed:     speed,
		battery:   battery,
		state:     state,
	}

	return dev, nil
}

// NewTracker creates a new GPS tracking device.
func NewTracker() *Device {
	dev, _ := NewDevice(DefaultTrackerOptions())
	return dev
}

// NewKidsTracker creates a new GPS tracking device.
func NewKidsTracker() *Device {
	dev, _ := NewDevice(KidsTrackerOptions())
	return dev
}

// NewDogTracker creates a new GPS tracking device.
func NewDogTracker() *Device {
	dev, _ := NewDevice(DogTrackerOptions())
	return dev
}

// NewBicycleTracker creates a new GPS tracking device.
func NewBicycleTracker() *Device {
	dev, _ := NewDevice(BicycleTrackerOptions())
	return dev
}

// NewDroneTracker creates a new GPS tracking device.
func NewDroneTracker() *Device {
	dev, _ := NewDevice(DroneTrackerOptions())
	return dev
}

// SetModel sets the model of the device.
func (d *Device) SetModel(model string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, err := types.NewModel(model)
	if err != nil {
		return err
	}
	d.state.Model = m.String()
	return nil
}

// SetModel sets the model of the device.
func (d *Device) SetUserID(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state.UserId = id
	return nil
}

// SetDescription sets the description of the device.
func (d *Device) SetDescription(descr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state.Description = descr
}

// ID returns the ID of the device.
func (d *Device) ID() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.Id
}

// Model returns the model of the device.
func (d *Device) Model() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.Model
}

// UserID returns the user ID associated with the device.
func (d *Device) UserID() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.UserId
}

// SetColor sets the color of the device.
func (d *Device) SetColor(color colorful.Color) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !color.IsValid() {
		return fmt.Errorf("gpsgen: invalid device color")
	}
	d.state.Color = color.Hex()
	return nil
}

// Color returns the color of the device.
func (d *Device) Color() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.Color
}

func (d *Device) Descr() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.Description
}

// AddSensor adds a sensor to the device with the specified parameters.
func (d *Device) AddSensor(
	name string,
	min, max float64,
	amplitude int,
	mode types.SensorMode,
) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sensor, err := types.NewSensor(name, min, max, amplitude, mode)
	if err != nil {
		return "", err
	}
	d.sensors = append(d.sensors, sensor)
	d.updateSensors()

	return sensor.ID(), nil
}

// SensorByID returns the sensor with the given ID
// and a boolean indicating its existence.
func (d *Device) SensorByID(sensorID string) (*types.Sensor, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for i := 0; i < len(d.sensors); i++ {
		if d.sensors[i].ID() == sensorID {
			return d.sensors[i], true
		}
	}
	return nil, false
}

// RemoveSensor removes a sensor with the specified sensorID from the device's list of sensors.
// It returns true if the sensor was successfully removed, and false otherwise.
// If the provided sensorID is empty, the function returns false.
func (d *Device) RemoveSensor(sensorID string) (ok bool) {
	if len(sensorID) == 0 {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// Search for the sensor by its ID in the device's list of sensors.
	for i := 0; i < len(d.sensors); i++ {
		if d.sensors[i].ID() == sensorID {
			ok = true
			// Remove the sensor from the slice by creating a new slice
			// that excludes the sensor being removed.
			d.sensors = append(d.sensors[:i], d.sensors[i+1:]...)
			break
		}
	}

	// If the sensor was successfully removed, update the device's sensors
	// and check if the sensors list is empty, in which case reset it.
	if ok {
		d.updateSensors()
		if len(d.sensors) == 0 {
			d.sensors = make([]*types.Sensor, 0)
		}
	}

	return
}

// NumRoutes returns the number of routes stored in the device's navigator.
func (d *Device) NumRoutes() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.NumRoutes()
}

// NumSensors returns the number of sensors attached to the device.
func (d *Device) NumSensors() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.sensors)
}

// SensorAt returns the sensor at the specified index from the device's navigator.
func (d *Device) SensorAt(i int) *types.Sensor {
	if len(d.sensors) == 0 || i > len(d.sensors)-1 || i < 0 {
		return nil
	}
	return d.sensors[i]
}

// EachSensor iterates over each sensor in the device's navigator and applies a function.
func (d *Device) EachSensor(fn func(int, *types.Sensor) bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for i := 0; i < len(d.sensors); i++ {
		if ok := fn(i, d.sensors[i]); !ok {
			break
		}
	}
}

// EachRoute iterates over each route in the device's navigator and applies a function.
func (d *Device) EachRoute(fn func(int, *navigator.Route) bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	d.navigator.EachRoute(fn)
}

// AddRoute adds one or more routes to the device's navigator.
func (d *Device) AddRoute(routes ...*navigator.Route) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	err := d.navigator.AddRoute(routes...)
	if err == nil {
		d.updateRoutes()
	}
	return err
}

// RemoveRoute removes a route from the device's navigator by its ID.
func (d *Device) RemoveRoute(routeID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	ok := d.navigator.RemoveRoute(routeID)
	if ok {
		d.updateRoutes()
	}
	return ok
}

// RemoveTrack removes a track from a route in the device's navigator.
func (d *Device) RemoveTrack(routeID, trackID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	ok := d.navigator.RemoveTrack(routeID, trackID)
	if ok {
		d.updateRoutes()
	}
	return ok
}

// Status returns the current status of the device.
func (d *Device) Status() Status {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.status
}

// State returns a protobuf representation of the current device state.
func (d *Device) State() *pb.Device {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

// Update updates the device's state and routes based on its current navigator state.
func (d *Device) Update() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.updateState()
	d.updateRoutes()
	d.updateSensors()
}

// CurrentBearing returns the current bearing direction of the device.
func (d *Device) CurrentBearing() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentBearing()
}

// Distance returns the total distance traveled by the device.
func (d *Device) Distance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.Distance()
}

// CurrentDistance returns the distance traveled in the current segment.
func (d *Device) CurrentDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentDistance()
}

// RouteDistance returns the total distance of the active route.
func (d *Device) RouteDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.RouteDistance()
}

// CurrentRouteDistance returns the distance traveled in the current route.
func (d *Device) CurrentRouteDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentRouteDistance()
}

// TrackDistance returns the total distance of the active track.
func (d *Device) TrackDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.TrackDistance()
}

// CurrentTrackDistance returns the distance traveled in the current track.
func (d *Device) CurrentTrackDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentTrackDistance()
}

// SegmentDistance returns the total distance of the active segment.
func (d *Device) SegmentDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.SegmentDistance()
}

// CurrentSegmentDistance returns the distance traveled in the current segment.
func (d *Device) CurrentSegmentDistance() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentSegmentDistance()
}

// IsFinish checks if the device has reached the end of its current route.
func (d *Device) IsFinish() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.IsFinish()
}

// CurrentRoute returns the currently active route of the device.
func (d *Device) CurrentRoute() *navigator.Route {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentRoute()
}

// CurrentTrack returns the currently active track of the device.
func (d *Device) CurrentTrack() *navigator.Track {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentTrack()
}

// CurrentSegment returns the currently active segment of the device.
func (d *Device) CurrentSegment() navigator.Segment {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.CurrentSegment()
}

// RouteIndex returns the index of the currently active route in the navigator.
func (d *Device) RouteIndex() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.RouteIndex()
}

// TrackIndex returns the index of the currently active track in the navigator.
func (d *Device) TrackIndex() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.TrackIndex()
}

// SegmentIndex returns the index of the current segment within the current track of the navigator.
func (d *Device) SegmentIndex() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.SegmentIndex()
}

// IsOffline returns whether the device is currently offline based on the navigator's status.
func (d *Device) IsOffline() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.IsOffline()
}

// ResetRoutes resets the routes associated with the device in the navigator.
// It returns true on success and updates the device's internal routes.
func (d *Device) ResetRoutes() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	ok := d.navigator.ResetRoutes()
	if ok {
		d.updateRoutes()
	}
	return ok
}

// Location returns the current geographical location of the device.
func (d *Device) Location() geo.LatLonPoint {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.Location()
}

// DestinationTo updates the device's position to a specified distance along the current segment.
// It returns true on success, indicating the update was applied.
func (d *Device) DestinationTo(meters float64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.DestinationTo(meters)
}

// MoveToSegment updates the device's position to a specific segment within a track and route.
// It returns true on success, indicating the update was applied.
func (d *Device) MoveToSegment(routeIndex int, trackIndex int, segmentIndex int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.MoveToSegment(routeIndex, trackIndex, segmentIndex)
}

// MoveToTrack updates the device's position to a specific track within a route.
// It returns true on success, indicating the update was applied.
func (d *Device) MoveToTrack(routeIndex int, trackIndex int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.MoveToTrack(routeIndex, trackIndex)
}

// MoveToRoute updates the device's position to a specific route.
// It returns true on success, indicating the update was applied.
func (d *Device) MoveToRoute(routeIndex int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.MoveToRoute(routeIndex)
}

// MoveToTrackByID updates the device's position to a specific track within a route, identified by IDs.
// It returns true on success, indicating the update was applied.
func (d *Device) MoveToTrackByID(routeID string, trackID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.MoveToTrackByID(routeID, trackID)
}

// MoveToRouteByID updates the device's position to a specific route, identified by its ID.
// It returns true on success, indicating the update was applied.
func (d *Device) MoveToRouteByID(routeID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.MoveToRouteByID(routeID)
}

// ToPrevRoute updates the device's position to the previous route in the navigator.
// It returns true on success, indicating the update was applied.
func (d *Device) ToPrevRoute() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.PrevRoute()
}

// ToNextRoute updates the device's position to the next route in the navigator.
// It returns true on success, indicating the update was applied.
func (d *Device) ToNextRoute() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.navigator.NextRoute()
}

// ToOffline sets the device's status to offline in the navigator.
// This method does not return a value.
func (d *Device) ToOffline() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.navigator.ToOffline()
}

// ResetNavigator resets the navigator for the device.
// This method clears all navigation-related state and resets the device's position.
// It does not return a value.
func (d *Device) ResetNavigator() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.navigator.Reset()
}

// RouteAt returns the route at the specified index in the navigator.
// The index parameter is the index of the desired route.
// It returns the route at the given index or nil if the index is out of bounds.
func (d *Device) RouteAt(index int) *navigator.Route {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.RouteAt(index)
}

// Duration returns the duration of the device's current track segment in seconds.
// It returns the duration of the current track segment or 0 if no track is active.
func (d *Device) Duration() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state.Duration
}

// Version returns the current version of the device as an array of three integers.
// It retrieves the version from the device's navigator and returns it as [Navigator, Route, Track].
func (d *Device) Version() [3]int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.navigator.Sum()
}

// Next advances the device's state for the next time step.
// The tick parameter represents the time interval in seconds since the last update.
// It updates the device's position, speed, elevation, battery level, sensors, and state.
// Returns true if the state was successfully updated, and false if no routes are available or if the navigator has finished.
func (d *Device) Next(tick float64) bool {
	ok := d.mu.TryLock()
	if !ok {
		d.incTick()
		return false
	}
	defer d.mu.Unlock()

	if d.navigator.NumRoutes() == 0 {
		return false
	}

	tickLock := atomic.LoadUint32(&d.tick)
	seconds := float64(float64(tickLock) + tick)
	d.resetTick()
	d.state.Duration += seconds
	d.state.Tick = seconds

	d.battery.Next(seconds)
	if d.battery.IsLow() {
		d.battery.Reset()
		d.navigator.ToOffline()
	}

	t := d.nextT(seconds)
	d.speed.Next(t)
	nextSpeed := d.speed.Value()
	d.state.Speed = nextSpeed

	if ok := d.navigator.NextLocation(seconds, nextSpeed); !ok {
		d.updateOfflineState()
		return false
	}

	d.navigator.NextElevation(t)

	if len(d.sensors) > 0 {
		for i := 0; i < len(d.sensors); i++ {
			d.sensors[i].Next(t)
		}
	}

	ns := d.navigator.Sum()
	if d.isNotValidNS(ns) {
		d.updateRoutes()
	}

	d.updateState()

	if d.navigator.IsFinish() {
		d.state.Duration = 0
		d.speed.Shuffle()
		for i := 0; i < len(d.sensors); i++ {
			d.sensors[i].Shuffle()
		}
	}

	return true
}

// MarshalBinary converts the device's current state into a binary representation.
// It creates a snapshot of the device and serializes it using protobuf encoding.
// Returns the binary representation of the device's state and an error, if any.
func (d *Device) MarshalBinary() ([]byte, error) {
	snap := d.Snapshot()
	return proto.Marshal(snap)
}

// UnmarshalBinary populates the device's state using the provided binary data.
// It deserializes the binary data using protobuf decoding and updates the device's state accordingly.
// Returns an error if the binary data is invalid or if there's an issue during deserialization.
func (d *Device) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("gpsgen: invalid snapshot data")
	}
	snap := new(pb.Snapshot)
	if err := proto.Unmarshal(data, snap); err != nil {
		return err
	}
	d.FromSnapshot(snap)
	return nil
}

// Snapshot creates and returns a snapshot of the device's current state.
// The snapshot includes device information, status, duration, navigator, speed, battery, and sensors (if any).
func (d *Device) Snapshot() *pb.Snapshot {
	snap := &pb.Snapshot{
		Id:        d.state.Id,
		UserId:    d.state.UserId,
		Model:     d.state.Model,
		Tick:      d.state.Tick,
		Color:     d.state.Color,
		Descr:     d.state.Description,
		Status:    int64(d.status),
		Duration:  d.state.Duration,
		Navigator: d.navigator.Snapshot(),
		Speed:     d.speed.Snapshot(),
		Battery:   d.battery.Snapshot(),
	}
	if len(d.sensors) > 0 {
		snap.Sensors = make([]*pb.Snapshot_SensorType, len(d.sensors))
		d.mu.RLock()
		for i := 0; i < len(d.sensors); i++ {
			snap.Sensors[i] = d.sensors[i].Snapshot()
		}
		d.mu.RUnlock()
	}
	return snap
}

// FromSnapshot updates the device's state using the provided snapshot data.
// It populates the device's state with the values
// from the snapshot, including device information, status, navigator, speed, battery, and sensors.
func (d *Device) FromSnapshot(snap *pb.Snapshot) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = initDeviceState()
	d.state.Id = snap.Id
	d.state.UserId = snap.UserId
	d.state.Tick = snap.Tick
	d.state.Model = snap.Model
	d.state.Color = snap.Color
	d.state.Description = snap.Descr
	d.state.Duration = snap.Duration
	d.status = Status(snap.Status)
	d.navigator = new(navigator.Navigator)
	d.navigator.FromSnapshot(snap.Navigator)
	d.speed = new(types.Speed)
	d.speed.FromSnapshot(snap.Speed)
	d.state.Speed = d.speed.Value()
	d.battery = new(types.Battery)
	d.battery.FromSnapshot(snap.Battery)
	if len(snap.Sensors) > 0 {
		d.sensors = make([]*types.Sensor, len(snap.Sensors))
		for i := 0; i < len(snap.Sensors); i++ {
			sensor := new(types.Sensor)
			sensor.FromSnapshot(snap.Sensors[i])
			d.sensors[i] = sensor
		}
	}
	d.updateSensors()
	d.updateState()
	d.updateRoutes()
}

func (d *Device) mount() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == Running {
		return fmt.Errorf("gpsgen: device %s is already running", d.state.Id)
	}
	d.status = Running
	return nil
}

func (d *Device) unmount() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == Stopped {
		return fmt.Errorf("gpsgen: device %s is already stopped", d.state.Id)
	}
	d.status = Stopped
	return nil
}

func (d *Device) incTick() {
	atomic.AddUint32(&d.tick, 1)
}

func (d *Device) resetTick() {
	atomic.StoreUint32(&d.tick, 0)
}

func (d *Device) nextT(tick float64) float64 {
	dur := d.state.Duration + tick
	return math.Min(dur/d.avgT, 1.0)
}

func (d *Device) updateState() {
	d.navigator.Update(d.state)
	d.state.Battery.Charge = d.battery.Value()
	d.state.Battery.ChargeTime = int64(d.battery.ChargeTime().Seconds())
	d.state.TimeEstimate = d.state.Distance.Distance / d.state.Speed
	if len(d.sensors) != len(d.state.Sensors) {
		d.updateSensors()
	} else {
		for i := 0; i < len(d.sensors); i++ {
			sensor := d.sensors[i]
			d.state.Sensors[i].ValX = sensor.ValueX()
			d.state.Sensors[i].ValY = sensor.ValueY()
		}
	}
}

func (d *Device) updateOfflineState() {
	d.state.OfflineDuration = int64(d.navigator.OfflineDuration())
	d.state.IsOffline = d.navigator.IsOffline()
}

func (d *Device) updateSensors() {
	d.state.Sensors = make([]*pb.Device_Sensor, 0, len(d.sensors))
	for i := 0; i < len(d.sensors); i++ {
		sensor := d.sensors[i]
		d.state.Sensors = append(d.state.Sensors, &pb.Device_Sensor{
			Id:   sensor.ID(),
			Name: sensor.Name(),
			ValX: sensor.ValueX(),
			ValY: sensor.ValueY(),
		})
	}
}

func (d *Device) updateRoutes() {
	d.state.Routes.Routes = make([]*pb.Device_Routes_Route, 0, d.navigator.NumRoutes())
	for i := 0; i < d.navigator.NumRoutes(); i++ {
		route := d.navigator.RouteAt(i)
		var routeProps []byte
		if len(route.Props()) > 0 {
			routeProps, _ = json.Marshal(route.Props())
		}
		tracks := make([]*pb.Device_Routes_Route_Track, 0, route.NumTracks())
		for j := 0; j < route.NumTracks(); j++ {
			track := route.TrackAt(j)
			var trackProps []byte
			if len(track.Props()) > 0 {
				trackProps, _ = json.Marshal(track.Props())
			}
			tracks = append(tracks, &pb.Device_Routes_Route_Track{
				Distance:    track.Distance(),
				NumSegments: int64(track.NumSegments()),
				Color:       track.Color(),
				Props:       trackProps,
				PropsCount:  int64(len(track.Props())),
			})
		}
		d.state.Routes.Routes = append(d.state.Routes.Routes,
			&pb.Device_Routes_Route{
				Id:         route.ID(),
				Tracks:     tracks,
				Distance:   route.Distance(),
				Color:      route.Color(),
				Props:      routeProps,
				PropsCount: int64(len(route.Props())),
			})
	}
	d.calcAvgT()
	d.ns = d.navigator.Sum()
}

func (d *Device) calcAvgT() {
	avgSpeed := (d.speed.Min() + d.speed.Max()) / 2
	d.avgT = d.navigator.Distance() / avgSpeed
}

func (d *Device) isNotValidNS(sum [3]int) (ok bool) {
	for i := 0; i < 3; i++ {
		if d.ns[i] != sum[i] {
			ok = true
			break
		}
	}
	return
}

func initDeviceState() *pb.Device {
	return &pb.Device{
		Routes:    &pb.Device_Routes{},
		Distance:  &pb.Device_Distance{},
		Navigator: &pb.Device_Navigator{},
		Battery:   &pb.Device_Battery{},
		Location: &pb.Device_Location{
			LatDms: &pb.Device_Location_DMS{},
			LonDms: &pb.Device_Location_DMS{},
			Utm:    &pb.Device_Location_UTM{},
		},
	}
}
