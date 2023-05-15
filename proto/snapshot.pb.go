// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: proto/snapshot.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeviceSnapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        []byte          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Model     string          `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	Speed     *TypeState      `protobuf:"bytes,3,opt,name=speed,proto3" json:"speed,omitempty"`
	Battery   *TypeState      `protobuf:"bytes,4,opt,name=battery,proto3" json:"battery,omitempty"`
	Sensors   []*SensorState  `protobuf:"bytes,5,rep,name=sensors,proto3" json:"sensors,omitempty"`
	Navigator *NavigatorState `protobuf:"bytes,6,opt,name=navigator,proto3" json:"navigator,omitempty"`
	Loop      float64         `protobuf:"fixed64,7,opt,name=loop,proto3" json:"loop,omitempty"`
	AvgTick   float64         `protobuf:"fixed64,8,opt,name=avg_tick,json=avgTick,proto3" json:"avg_tick,omitempty"`
}

func (x *DeviceSnapshot) Reset() {
	*x = DeviceSnapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceSnapshot) ProtoMessage() {}

func (x *DeviceSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceSnapshot.ProtoReflect.Descriptor instead.
func (*DeviceSnapshot) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceSnapshot) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DeviceSnapshot) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *DeviceSnapshot) GetSpeed() *TypeState {
	if x != nil {
		return x.Speed
	}
	return nil
}

func (x *DeviceSnapshot) GetBattery() *TypeState {
	if x != nil {
		return x.Battery
	}
	return nil
}

func (x *DeviceSnapshot) GetSensors() []*SensorState {
	if x != nil {
		return x.Sensors
	}
	return nil
}

func (x *DeviceSnapshot) GetNavigator() *NavigatorState {
	if x != nil {
		return x.Navigator
	}
	return nil
}

func (x *DeviceSnapshot) GetLoop() float64 {
	if x != nil {
		return x.Loop
	}
	return 0
}

func (x *DeviceSnapshot) GetAvgTick() float64 {
	if x != nil {
		return x.AvgTick
	}
	return 0
}

type NavigatorState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Routes          []*NavigatorState_Route `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
	RouteIndex      int64                   `protobuf:"varint,2,opt,name=route_index,json=routeIndex,proto3" json:"route_index,omitempty"`
	TrackIndex      int64                   `protobuf:"varint,3,opt,name=track_index,json=trackIndex,proto3" json:"track_index,omitempty"`
	SegmentIndex    int64                   `protobuf:"varint,4,opt,name=segment_index,json=segmentIndex,proto3" json:"segment_index,omitempty"`
	SegmentDistance float64                 `protobuf:"fixed64,5,opt,name=segment_distance,json=segmentDistance,proto3" json:"segment_distance,omitempty"`
	CurrentDistance float64                 `protobuf:"fixed64,6,opt,name=current_distance,json=currentDistance,proto3" json:"current_distance,omitempty"`
	OfflineIndex    int64                   `protobuf:"varint,7,opt,name=offline_index,json=offlineIndex,proto3" json:"offline_index,omitempty"`
	Point           *NavigatorState_Point   `protobuf:"bytes,8,opt,name=point,proto3" json:"point,omitempty"`
	Elevation       *SensorState            `protobuf:"bytes,9,opt,name=elevation,proto3" json:"elevation,omitempty"`
	OfflineMin      int64                   `protobuf:"varint,10,opt,name=offline_min,json=offlineMin,proto3" json:"offline_min,omitempty"`
	OfflineMax      int64                   `protobuf:"varint,11,opt,name=offline_max,json=offlineMax,proto3" json:"offline_max,omitempty"`
	TotalDistance   float64                 `protobuf:"fixed64,12,opt,name=total_distance,json=totalDistance,proto3" json:"total_distance,omitempty"`
	SkipOffline     bool                    `protobuf:"varint,13,opt,name=skip_offline,json=skipOffline,proto3" json:"skip_offline,omitempty"`
}

func (x *NavigatorState) Reset() {
	*x = NavigatorState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorState) ProtoMessage() {}

func (x *NavigatorState) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorState.ProtoReflect.Descriptor instead.
func (*NavigatorState) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{1}
}

func (x *NavigatorState) GetRoutes() []*NavigatorState_Route {
	if x != nil {
		return x.Routes
	}
	return nil
}

func (x *NavigatorState) GetRouteIndex() int64 {
	if x != nil {
		return x.RouteIndex
	}
	return 0
}

func (x *NavigatorState) GetTrackIndex() int64 {
	if x != nil {
		return x.TrackIndex
	}
	return 0
}

func (x *NavigatorState) GetSegmentIndex() int64 {
	if x != nil {
		return x.SegmentIndex
	}
	return 0
}

func (x *NavigatorState) GetSegmentDistance() float64 {
	if x != nil {
		return x.SegmentDistance
	}
	return 0
}

func (x *NavigatorState) GetCurrentDistance() float64 {
	if x != nil {
		return x.CurrentDistance
	}
	return 0
}

func (x *NavigatorState) GetOfflineIndex() int64 {
	if x != nil {
		return x.OfflineIndex
	}
	return 0
}

func (x *NavigatorState) GetPoint() *NavigatorState_Point {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *NavigatorState) GetElevation() *SensorState {
	if x != nil {
		return x.Elevation
	}
	return nil
}

func (x *NavigatorState) GetOfflineMin() int64 {
	if x != nil {
		return x.OfflineMin
	}
	return 0
}

func (x *NavigatorState) GetOfflineMax() int64 {
	if x != nil {
		return x.OfflineMax
	}
	return 0
}

func (x *NavigatorState) GetTotalDistance() float64 {
	if x != nil {
		return x.TotalDistance
	}
	return 0
}

func (x *NavigatorState) GetSkipOffline() bool {
	if x != nil {
		return x.SkipOffline
	}
	return false
}

type TypeState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min float64 `protobuf:"fixed64,1,opt,name=min,proto3" json:"min,omitempty"`
	Max float64 `protobuf:"fixed64,2,opt,name=max,proto3" json:"max,omitempty"`
	Val float64 `protobuf:"fixed64,3,opt,name=val,proto3" json:"val,omitempty"`
	Gen *Curve  `protobuf:"bytes,4,opt,name=gen,proto3" json:"gen,omitempty"`
}

func (x *TypeState) Reset() {
	*x = TypeState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypeState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypeState) ProtoMessage() {}

func (x *TypeState) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypeState.ProtoReflect.Descriptor instead.
func (*TypeState) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{2}
}

func (x *TypeState) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *TypeState) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *TypeState) GetVal() float64 {
	if x != nil {
		return x.Val
	}
	return 0
}

func (x *TypeState) GetGen() *Curve {
	if x != nil {
		return x.Gen
	}
	return nil
}

type SensorState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min  float64 `protobuf:"fixed64,1,opt,name=min,proto3" json:"min,omitempty"`
	Max  float64 `protobuf:"fixed64,2,opt,name=max,proto3" json:"max,omitempty"`
	ValX float64 `protobuf:"fixed64,3,opt,name=val_x,json=valX,proto3" json:"val_x,omitempty"`
	ValY float64 `protobuf:"fixed64,4,opt,name=val_y,json=valY,proto3" json:"val_y,omitempty"`
	Name string  `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Gen  *Curve  `protobuf:"bytes,6,opt,name=gen,proto3" json:"gen,omitempty"`
}

func (x *SensorState) Reset() {
	*x = SensorState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorState) ProtoMessage() {}

func (x *SensorState) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorState.ProtoReflect.Descriptor instead.
func (*SensorState) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{3}
}

func (x *SensorState) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *SensorState) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *SensorState) GetValX() float64 {
	if x != nil {
		return x.ValX
	}
	return 0
}

func (x *SensorState) GetValY() float64 {
	if x != nil {
		return x.ValY
	}
	return 0
}

func (x *SensorState) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SensorState) GetGen() *Curve {
	if x != nil {
		return x.Gen
	}
	return nil
}

type Curve struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Points []*Curve_ControlPoint `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
	Mode   int64                 `protobuf:"varint,2,opt,name=mode,proto3" json:"mode,omitempty"`
}

func (x *Curve) Reset() {
	*x = Curve{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Curve) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Curve) ProtoMessage() {}

func (x *Curve) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Curve.ProtoReflect.Descriptor instead.
func (*Curve) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{4}
}

func (x *Curve) GetPoints() []*Curve_ControlPoint {
	if x != nil {
		return x.Points
	}
	return nil
}

func (x *Curve) GetMode() int64 {
	if x != nil {
		return x.Mode
	}
	return 0
}

type NavigatorState_Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lat float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon float64 `protobuf:"fixed64,2,opt,name=lon,proto3" json:"lon,omitempty"`
}

func (x *NavigatorState_Point) Reset() {
	*x = NavigatorState_Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorState_Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorState_Point) ProtoMessage() {}

func (x *NavigatorState_Point) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorState_Point.ProtoReflect.Descriptor instead.
func (*NavigatorState_Point) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{1, 0}
}

func (x *NavigatorState_Point) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *NavigatorState_Point) GetLon() float64 {
	if x != nil {
		return x.Lon
	}
	return 0
}

type NavigatorState_Route struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Distance float64                       `protobuf:"fixed64,1,opt,name=distance,proto3" json:"distance,omitempty"`
	Tracks   []*NavigatorState_Route_Track `protobuf:"bytes,2,rep,name=tracks,proto3" json:"tracks,omitempty"`
}

func (x *NavigatorState_Route) Reset() {
	*x = NavigatorState_Route{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorState_Route) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorState_Route) ProtoMessage() {}

func (x *NavigatorState_Route) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorState_Route.ProtoReflect.Descriptor instead.
func (*NavigatorState_Route) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{1, 1}
}

func (x *NavigatorState_Route) GetDistance() float64 {
	if x != nil {
		return x.Distance
	}
	return 0
}

func (x *NavigatorState_Route) GetTracks() []*NavigatorState_Route_Track {
	if x != nil {
		return x.Tracks
	}
	return nil
}

type NavigatorState_Route_Track struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Segmenets []*NavigatorState_Route_Track_Segment `protobuf:"bytes,1,rep,name=segmenets,proto3" json:"segmenets,omitempty"`
}

func (x *NavigatorState_Route_Track) Reset() {
	*x = NavigatorState_Route_Track{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorState_Route_Track) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorState_Route_Track) ProtoMessage() {}

func (x *NavigatorState_Route_Track) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorState_Route_Track.ProtoReflect.Descriptor instead.
func (*NavigatorState_Route_Track) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{1, 1, 0}
}

func (x *NavigatorState_Route_Track) GetSegmenets() []*NavigatorState_Route_Track_Segment {
	if x != nil {
		return x.Segmenets
	}
	return nil
}

type NavigatorState_Route_Track_Segment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PointA   *NavigatorState_Point `protobuf:"bytes,1,opt,name=point_a,json=pointA,proto3" json:"point_a,omitempty"`
	PointB   *NavigatorState_Point `protobuf:"bytes,2,opt,name=point_b,json=pointB,proto3" json:"point_b,omitempty"`
	Distance float64               `protobuf:"fixed64,3,opt,name=distance,proto3" json:"distance,omitempty"`
	Bearing  float64               `protobuf:"fixed64,4,opt,name=bearing,proto3" json:"bearing,omitempty"`
	Rel      int64                 `protobuf:"varint,5,opt,name=rel,proto3" json:"rel,omitempty"`
}

func (x *NavigatorState_Route_Track_Segment) Reset() {
	*x = NavigatorState_Route_Track_Segment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorState_Route_Track_Segment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorState_Route_Track_Segment) ProtoMessage() {}

func (x *NavigatorState_Route_Track_Segment) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorState_Route_Track_Segment.ProtoReflect.Descriptor instead.
func (*NavigatorState_Route_Track_Segment) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{1, 1, 0, 0}
}

func (x *NavigatorState_Route_Track_Segment) GetPointA() *NavigatorState_Point {
	if x != nil {
		return x.PointA
	}
	return nil
}

func (x *NavigatorState_Route_Track_Segment) GetPointB() *NavigatorState_Point {
	if x != nil {
		return x.PointB
	}
	return nil
}

func (x *NavigatorState_Route_Track_Segment) GetDistance() float64 {
	if x != nil {
		return x.Distance
	}
	return 0
}

func (x *NavigatorState_Route_Track_Segment) GetBearing() float64 {
	if x != nil {
		return x.Bearing
	}
	return 0
}

func (x *NavigatorState_Route_Track_Segment) GetRel() int64 {
	if x != nil {
		return x.Rel
	}
	return 0
}

type Curve_Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Curve_Point) Reset() {
	*x = Curve_Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Curve_Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Curve_Point) ProtoMessage() {}

func (x *Curve_Point) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Curve_Point.ProtoReflect.Descriptor instead.
func (*Curve_Point) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{4, 0}
}

func (x *Curve_Point) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Curve_Point) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

type Curve_ControlPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vp *Curve_Point `protobuf:"bytes,1,opt,name=vp,proto3" json:"vp,omitempty"`
	Cp *Curve_Point `protobuf:"bytes,2,opt,name=cp,proto3" json:"cp,omitempty"`
}

func (x *Curve_ControlPoint) Reset() {
	*x = Curve_ControlPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_snapshot_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Curve_ControlPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Curve_ControlPoint) ProtoMessage() {}

func (x *Curve_ControlPoint) ProtoReflect() protoreflect.Message {
	mi := &file_proto_snapshot_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Curve_ControlPoint.ProtoReflect.Descriptor instead.
func (*Curve_ControlPoint) Descriptor() ([]byte, []int) {
	return file_proto_snapshot_proto_rawDescGZIP(), []int{4, 1}
}

func (x *Curve_ControlPoint) GetVp() *Curve_Point {
	if x != nil {
		return x.Vp
	}
	return nil
}

func (x *Curve_ControlPoint) GetCp() *Curve_Point {
	if x != nil {
		return x.Cp
	}
	return nil
}

var File_proto_snapshot_proto protoreflect.FileDescriptor

var file_proto_snapshot_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x02,
	0x0a, 0x0e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x26, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x12, 0x2a,
	0x0a, 0x07, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x07, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x12, 0x2c, 0x0a, 0x07, 0x73, 0x65,
	0x6e, 0x73, 0x6f, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x07, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x12, 0x33, 0x0a, 0x09, 0x6e, 0x61, 0x76, 0x69,
	0x67, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x09, 0x6e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x6f, 0x6f, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x6f, 0x6f,
	0x70, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x76, 0x67, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x76, 0x67, 0x54, 0x69, 0x63, 0x6b, 0x22, 0xb9, 0x07, 0x0a,
	0x0e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x33, 0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x06, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x29, 0x0a, 0x10, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x69,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e,
	0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61,
	0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x30, 0x0a, 0x09, 0x65, 0x6c, 0x65,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x09, 0x65, 0x6c, 0x65, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x6f,
	0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0a, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x4d, 0x69, 0x6e, 0x12, 0x1f, 0x0a, 0x0b,
	0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x4d, 0x61, 0x78, 0x12, 0x25, 0x0a,
	0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x6b, 0x69, 0x70, 0x5f, 0x6f, 0x66, 0x66,
	0x6c, 0x69, 0x6e, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x6b, 0x69, 0x70,
	0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x1a, 0x2b, 0x0a, 0x05, 0x50, 0x6f, 0x69, 0x6e, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c,
	0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x03, 0x6c, 0x6f, 0x6e, 0x1a, 0xf1, 0x02, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x74, 0x72,
	0x61, 0x63, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x06, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x73, 0x1a, 0x90, 0x02, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x12,
	0x47, 0x0a, 0x09, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67,
	0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x2e,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x65, 0x74, 0x73, 0x1a, 0xbd, 0x01, 0x0a, 0x07, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61,
	0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x41, 0x12, 0x34, 0x0a, 0x07, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x5f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x42,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x62, 0x65, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x62,
	0x65, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x6c, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x65, 0x6c, 0x22, 0x61, 0x0a, 0x09, 0x54, 0x79, 0x70, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x1e, 0x0a, 0x03, 0x67,
	0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x75, 0x72, 0x76, 0x65, 0x52, 0x03, 0x67, 0x65, 0x6e, 0x22, 0x8f, 0x01, 0x0a, 0x0b,
	0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d,
	0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12,
	0x13, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x5f, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04,
	0x76, 0x61, 0x6c, 0x58, 0x12, 0x13, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x5f, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x04, 0x76, 0x61, 0x6c, 0x59, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x03, 0x67, 0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x75, 0x72, 0x76, 0x65, 0x52, 0x03, 0x67, 0x65, 0x6e, 0x22, 0xcb, 0x01,
	0x0a, 0x05, 0x43, 0x75, 0x72, 0x76, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x43, 0x75, 0x72, 0x76, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x1a, 0x23,
	0x0a, 0x05, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x01, 0x79, 0x1a, 0x56, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x02, 0x76, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x75, 0x72, 0x76, 0x65, 0x2e, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x52, 0x02, 0x76, 0x70, 0x12, 0x22, 0x0a, 0x02, 0x63, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x75, 0x72, 0x76,
	0x65, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x02, 0x63, 0x70, 0x42, 0x09, 0x5a, 0x07, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_snapshot_proto_rawDescOnce sync.Once
	file_proto_snapshot_proto_rawDescData = file_proto_snapshot_proto_rawDesc
)

func file_proto_snapshot_proto_rawDescGZIP() []byte {
	file_proto_snapshot_proto_rawDescOnce.Do(func() {
		file_proto_snapshot_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_snapshot_proto_rawDescData)
	})
	return file_proto_snapshot_proto_rawDescData
}

var file_proto_snapshot_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_snapshot_proto_goTypes = []interface{}{
	(*DeviceSnapshot)(nil),                     // 0: proto.DeviceSnapshot
	(*NavigatorState)(nil),                     // 1: proto.NavigatorState
	(*TypeState)(nil),                          // 2: proto.TypeState
	(*SensorState)(nil),                        // 3: proto.SensorState
	(*Curve)(nil),                              // 4: proto.Curve
	(*NavigatorState_Point)(nil),               // 5: proto.NavigatorState.Point
	(*NavigatorState_Route)(nil),               // 6: proto.NavigatorState.Route
	(*NavigatorState_Route_Track)(nil),         // 7: proto.NavigatorState.Route.Track
	(*NavigatorState_Route_Track_Segment)(nil), // 8: proto.NavigatorState.Route.Track.Segment
	(*Curve_Point)(nil),                        // 9: proto.Curve.Point
	(*Curve_ControlPoint)(nil),                 // 10: proto.Curve.ControlPoint
}
var file_proto_snapshot_proto_depIdxs = []int32{
	2,  // 0: proto.DeviceSnapshot.speed:type_name -> proto.TypeState
	2,  // 1: proto.DeviceSnapshot.battery:type_name -> proto.TypeState
	3,  // 2: proto.DeviceSnapshot.sensors:type_name -> proto.SensorState
	1,  // 3: proto.DeviceSnapshot.navigator:type_name -> proto.NavigatorState
	6,  // 4: proto.NavigatorState.routes:type_name -> proto.NavigatorState.Route
	5,  // 5: proto.NavigatorState.point:type_name -> proto.NavigatorState.Point
	3,  // 6: proto.NavigatorState.elevation:type_name -> proto.SensorState
	4,  // 7: proto.TypeState.gen:type_name -> proto.Curve
	4,  // 8: proto.SensorState.gen:type_name -> proto.Curve
	10, // 9: proto.Curve.points:type_name -> proto.Curve.ControlPoint
	7,  // 10: proto.NavigatorState.Route.tracks:type_name -> proto.NavigatorState.Route.Track
	8,  // 11: proto.NavigatorState.Route.Track.segmenets:type_name -> proto.NavigatorState.Route.Track.Segment
	5,  // 12: proto.NavigatorState.Route.Track.Segment.point_a:type_name -> proto.NavigatorState.Point
	5,  // 13: proto.NavigatorState.Route.Track.Segment.point_b:type_name -> proto.NavigatorState.Point
	9,  // 14: proto.Curve.ControlPoint.vp:type_name -> proto.Curve.Point
	9,  // 15: proto.Curve.ControlPoint.cp:type_name -> proto.Curve.Point
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_proto_snapshot_proto_init() }
func file_proto_snapshot_proto_init() {
	if File_proto_snapshot_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_snapshot_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceSnapshot); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypeState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Curve); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorState_Point); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorState_Route); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorState_Route_Track); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorState_Route_Track_Segment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Curve_Point); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_snapshot_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Curve_ControlPoint); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_snapshot_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_snapshot_proto_goTypes,
		DependencyIndexes: file_proto_snapshot_proto_depIdxs,
		MessageInfos:      file_proto_snapshot_proto_msgTypes,
	}.Build()
	File_proto_snapshot_proto = out.File
	file_proto_snapshot_proto_rawDesc = nil
	file_proto_snapshot_proto_goTypes = nil
	file_proto_snapshot_proto_depIdxs = nil
}
