// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: common.proto

package game_proto

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

type Pos struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float64 `protobuf:"fixed64,1,opt,name=X,proto3" json:"X,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=Y,proto3" json:"Y,omitempty"`
}

func (x *Pos) Reset() {
	*x = Pos{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pos) ProtoMessage() {}

func (x *Pos) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pos.ProtoReflect.Descriptor instead.
func (*Pos) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *Pos) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Pos) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

// 坦克信息
type TankInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Level     int32   `protobuf:"varint,2,opt,name=Level,proto3" json:"Level,omitempty"`
	CurrPos   *Pos    `protobuf:"bytes,3,opt,name=CurrPos,proto3" json:"CurrPos,omitempty"`
	Direction int32   `protobuf:"varint,4,opt,name=Direction,proto3" json:"Direction,omitempty"`
	CurrSpeed float32 `protobuf:"fixed32,5,opt,name=CurrSpeed,proto3" json:"CurrSpeed,omitempty"`
}

func (x *TankInfo) Reset() {
	*x = TankInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TankInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TankInfo) ProtoMessage() {}

func (x *TankInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TankInfo.ProtoReflect.Descriptor instead.
func (*TankInfo) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *TankInfo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TankInfo) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *TankInfo) GetCurrPos() *Pos {
	if x != nil {
		return x.CurrPos
	}
	return nil
}

func (x *TankInfo) GetDirection() int32 {
	if x != nil {
		return x.Direction
	}
	return 0
}

func (x *TankInfo) GetCurrSpeed() float32 {
	if x != nil {
		return x.CurrSpeed
	}
	return 0
}

// 坦克移动信息
type TankMoveInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrPos              *Pos    `protobuf:"bytes,1,opt,name=CurrPos,proto3" json:"CurrPos,omitempty"`
	Direction            int32   `protobuf:"varint,2,opt,name=Direction,proto3" json:"Direction,omitempty"`
	CurrSpeed            float32 `protobuf:"fixed32,3,opt,name=CurrSpeed,proto3" json:"CurrSpeed,omitempty"`
	CurrTimeMilliseconds int64   `protobuf:"varint,4,opt,name=CurrTimeMilliseconds,proto3" json:"CurrTimeMilliseconds,omitempty"`
}

func (x *TankMoveInfo) Reset() {
	*x = TankMoveInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TankMoveInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TankMoveInfo) ProtoMessage() {}

func (x *TankMoveInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TankMoveInfo.ProtoReflect.Descriptor instead.
func (*TankMoveInfo) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

func (x *TankMoveInfo) GetCurrPos() *Pos {
	if x != nil {
		return x.CurrPos
	}
	return nil
}

func (x *TankMoveInfo) GetDirection() int32 {
	if x != nil {
		return x.Direction
	}
	return 0
}

func (x *TankMoveInfo) GetCurrSpeed() float32 {
	if x != nil {
		return x.CurrSpeed
	}
	return 0
}

func (x *TankMoveInfo) GetCurrTimeMilliseconds() int64 {
	if x != nil {
		return x.CurrTimeMilliseconds
	}
	return 0
}

// 玩家坦克信息
type PlayerTankInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=PlayerId,proto3" json:"PlayerId,omitempty"`
	TankInfo *TankInfo `protobuf:"bytes,2,opt,name=TankInfo,proto3" json:"TankInfo,omitempty"`
}

func (x *PlayerTankInfo) Reset() {
	*x = PlayerTankInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerTankInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerTankInfo) ProtoMessage() {}

func (x *PlayerTankInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerTankInfo.ProtoReflect.Descriptor instead.
func (*PlayerTankInfo) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

func (x *PlayerTankInfo) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *PlayerTankInfo) GetTankInfo() *TankInfo {
	if x != nil {
		return x.TankInfo
	}
	return nil
}

// 玩家帐号和坦克信息
type PlayerAccountTankInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=PlayerId,proto3" json:"PlayerId,omitempty"`
	Account  string    `protobuf:"bytes,2,opt,name=Account,proto3" json:"Account,omitempty"`
	TankInfo *TankInfo `protobuf:"bytes,3,opt,name=TankInfo,proto3" json:"TankInfo,omitempty"`
}

func (x *PlayerAccountTankInfo) Reset() {
	*x = PlayerAccountTankInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerAccountTankInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerAccountTankInfo) ProtoMessage() {}

func (x *PlayerAccountTankInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerAccountTankInfo.ProtoReflect.Descriptor instead.
func (*PlayerAccountTankInfo) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{4}
}

func (x *PlayerAccountTankInfo) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *PlayerAccountTankInfo) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *PlayerAccountTankInfo) GetTankInfo() *TankInfo {
	if x != nil {
		return x.TankInfo
	}
	return nil
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a, 0x03, 0x50, 0x6f,
	0x73, 0x12, 0x0c, 0x0a, 0x01, 0x58, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x58, 0x12,
	0x0c, 0x0a, 0x01, 0x59, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x59, 0x22, 0x97, 0x01,
	0x0a, 0x08, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c,
	0x12, 0x29, 0x0a, 0x07, 0x43, 0x75, 0x72, 0x72, 0x50, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x6f, 0x73, 0x52, 0x07, 0x43, 0x75, 0x72, 0x72, 0x50, 0x6f, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x44,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x75, 0x72,
	0x72, 0x53, 0x70, 0x65, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x43, 0x75,
	0x72, 0x72, 0x53, 0x70, 0x65, 0x65, 0x64, 0x22, 0xa9, 0x01, 0x0a, 0x0c, 0x54, 0x61, 0x6e, 0x6b,
	0x4d, 0x6f, 0x76, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x07, 0x43, 0x75, 0x72, 0x72,
	0x50, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x73, 0x52, 0x07, 0x43, 0x75, 0x72, 0x72,
	0x50, 0x6f, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x75, 0x72, 0x72, 0x53, 0x70, 0x65, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x43, 0x75, 0x72, 0x72, 0x53, 0x70, 0x65, 0x65, 0x64, 0x12,
	0x32, 0x0a, 0x14, 0x43, 0x75, 0x72, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x69, 0x6c, 0x6c, 0x69,
	0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x43,
	0x75, 0x72, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x22, 0x5e, 0x0a, 0x0e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x54, 0x61, 0x6e,
	0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x30, 0x0a, 0x08, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x54, 0x61, 0x6e, 0x6b, 0x49,
	0x6e, 0x66, 0x6f, 0x22, 0x7f, 0x0a, 0x15, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x61, 0x6e, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x54, 0x61, 0x6e, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_common_proto_goTypes = []interface{}{
	(*Pos)(nil),                   // 0: game_proto.Pos
	(*TankInfo)(nil),              // 1: game_proto.TankInfo
	(*TankMoveInfo)(nil),          // 2: game_proto.TankMoveInfo
	(*PlayerTankInfo)(nil),        // 3: game_proto.PlayerTankInfo
	(*PlayerAccountTankInfo)(nil), // 4: game_proto.PlayerAccountTankInfo
}
var file_common_proto_depIdxs = []int32{
	0, // 0: game_proto.TankInfo.CurrPos:type_name -> game_proto.Pos
	0, // 1: game_proto.TankMoveInfo.CurrPos:type_name -> game_proto.Pos
	1, // 2: game_proto.PlayerTankInfo.TankInfo:type_name -> game_proto.TankInfo
	1, // 3: game_proto.PlayerAccountTankInfo.TankInfo:type_name -> game_proto.TankInfo
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pos); i {
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
		file_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TankInfo); i {
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
		file_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TankMoveInfo); i {
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
		file_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerTankInfo); i {
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
		file_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerAccountTankInfo); i {
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
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
