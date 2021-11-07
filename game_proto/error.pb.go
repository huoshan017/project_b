// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: error.proto

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

// 错误ID定义
type ErrorId int32

const (
	ErrorId_NONE                                ErrorId = 0
	ErrorId_ACCOUNT_NOT_FOUND                   ErrorId = 1   // 帐号找不到
	ErrorId_INVALID_ACCOUNT                     ErrorId = 2   // 非法帐号名
	ErrorId_ACCOUNT_IS_LOGGIN                   ErrorId = 3   // 帐号正在登录
	ErrorId_SESSION_INTERNAL_ERROR              ErrorId = 100 // 会话内部错误
	ErrorId_PLAYER_ENTER_GAME_REPEATED          ErrorId = 101 // 玩家重复进入游戏
	ErrorId_DIFFERENT_PLAYER_ENTER_SAME_SESSION ErrorId = 102 // 不同玩家进入同一会话
	ErrorId_PLAYER_ENTERING_GAME                ErrorId = 103 // 玩家正在进入游戏
	ErrorId_PLAYER_CHANGE_TANK_FAILED           ErrorId = 200 // 玩家改变坦克失败
	ErrorId_PLAYER_RESTORE_TANK_FAILED          ErrorId = 201 // 玩家恢复坦克失败
)

// Enum value maps for ErrorId.
var (
	ErrorId_name = map[int32]string{
		0:   "NONE",
		1:   "ACCOUNT_NOT_FOUND",
		2:   "INVALID_ACCOUNT",
		3:   "ACCOUNT_IS_LOGGIN",
		100: "SESSION_INTERNAL_ERROR",
		101: "PLAYER_ENTER_GAME_REPEATED",
		102: "DIFFERENT_PLAYER_ENTER_SAME_SESSION",
		103: "PLAYER_ENTERING_GAME",
		200: "PLAYER_CHANGE_TANK_FAILED",
		201: "PLAYER_RESTORE_TANK_FAILED",
	}
	ErrorId_value = map[string]int32{
		"NONE":                                0,
		"ACCOUNT_NOT_FOUND":                   1,
		"INVALID_ACCOUNT":                     2,
		"ACCOUNT_IS_LOGGIN":                   3,
		"SESSION_INTERNAL_ERROR":              100,
		"PLAYER_ENTER_GAME_REPEATED":          101,
		"DIFFERENT_PLAYER_ENTER_SAME_SESSION": 102,
		"PLAYER_ENTERING_GAME":                103,
		"PLAYER_CHANGE_TANK_FAILED":           200,
		"PLAYER_RESTORE_TANK_FAILED":          201,
	}
)

func (x ErrorId) Enum() *ErrorId {
	p := new(ErrorId)
	*p = x
	return p
}

func (x ErrorId) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorId) Descriptor() protoreflect.EnumDescriptor {
	return file_error_proto_enumTypes[0].Descriptor()
}

func (ErrorId) Type() protoreflect.EnumType {
	return &file_error_proto_enumTypes[0]
}

func (x ErrorId) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorId.Descriptor instead.
func (ErrorId) EnumDescriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

type MsgErrorAck_ProtoId int32

const (
	MsgErrorAck_None MsgErrorAck_ProtoId = 0
	MsgErrorAck_Id   MsgErrorAck_ProtoId = 9999
)

// Enum value maps for MsgErrorAck_ProtoId.
var (
	MsgErrorAck_ProtoId_name = map[int32]string{
		0:    "None",
		9999: "Id",
	}
	MsgErrorAck_ProtoId_value = map[string]int32{
		"None": 0,
		"Id":   9999,
	}
)

func (x MsgErrorAck_ProtoId) Enum() *MsgErrorAck_ProtoId {
	p := new(MsgErrorAck_ProtoId)
	*p = x
	return p
}

func (x MsgErrorAck_ProtoId) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgErrorAck_ProtoId) Descriptor() protoreflect.EnumDescriptor {
	return file_error_proto_enumTypes[1].Descriptor()
}

func (MsgErrorAck_ProtoId) Type() protoreflect.EnumType {
	return &file_error_proto_enumTypes[1]
}

func (x MsgErrorAck_ProtoId) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgErrorAck_ProtoId.Descriptor instead.
func (MsgErrorAck_ProtoId) EnumDescriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0, 0}
}

// 错误协议返回
type MsgErrorAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error ErrorId `protobuf:"varint,1,opt,name=Error,proto3,enum=game_proto.ErrorId" json:"Error,omitempty"`
}

func (x *MsgErrorAck) Reset() {
	*x = MsgErrorAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgErrorAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgErrorAck) ProtoMessage() {}

func (x *MsgErrorAck) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgErrorAck.ProtoReflect.Descriptor instead.
func (*MsgErrorAck) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

func (x *MsgErrorAck) GetError() ErrorId {
	if x != nil {
		return x.Error
	}
	return ErrorId_NONE
}

var File_error_proto protoreflect.FileDescriptor

var file_error_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67,
	0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x0b, 0x4d, 0x73, 0x67,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x41, 0x63, 0x6b, 0x12, 0x29, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x64, 0x52, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x22, 0x1c, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x49, 0x64, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x02, 0x49, 0x64, 0x10, 0x8f,
	0x4e, 0x2a, 0x96, 0x02, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x08, 0x0a,
	0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x43, 0x4f, 0x55,
	0x4e, 0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x13,
	0x0a, 0x0f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x49,
	0x53, 0x5f, 0x4c, 0x4f, 0x47, 0x47, 0x49, 0x4e, 0x10, 0x03, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45,
	0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x64, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52,
	0x5f, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x52, 0x45, 0x50, 0x45,
	0x41, 0x54, 0x45, 0x44, 0x10, 0x65, 0x12, 0x27, 0x0a, 0x23, 0x44, 0x49, 0x46, 0x46, 0x45, 0x52,
	0x45, 0x4e, 0x54, 0x5f, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52, 0x5f, 0x45, 0x4e, 0x54, 0x45, 0x52,
	0x5f, 0x53, 0x41, 0x4d, 0x45, 0x5f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x66, 0x12,
	0x18, 0x0a, 0x14, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52, 0x5f, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x49,
	0x4e, 0x47, 0x5f, 0x47, 0x41, 0x4d, 0x45, 0x10, 0x67, 0x12, 0x1e, 0x0a, 0x19, 0x50, 0x4c, 0x41,
	0x59, 0x45, 0x52, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f, 0x54, 0x41, 0x4e, 0x4b, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0xc8, 0x01, 0x12, 0x1f, 0x0a, 0x1a, 0x50, 0x4c, 0x41,
	0x59, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x54, 0x4f, 0x52, 0x45, 0x5f, 0x54, 0x41, 0x4e, 0x4b,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0xc9, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2e,
	0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_error_proto_rawDescOnce sync.Once
	file_error_proto_rawDescData = file_error_proto_rawDesc
)

func file_error_proto_rawDescGZIP() []byte {
	file_error_proto_rawDescOnce.Do(func() {
		file_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_error_proto_rawDescData)
	})
	return file_error_proto_rawDescData
}

var file_error_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_error_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_error_proto_goTypes = []interface{}{
	(ErrorId)(0),             // 0: game_proto.ErrorId
	(MsgErrorAck_ProtoId)(0), // 1: game_proto.MsgErrorAck.ProtoId
	(*MsgErrorAck)(nil),      // 2: game_proto.MsgErrorAck
}
var file_error_proto_depIdxs = []int32{
	0, // 0: game_proto.MsgErrorAck.Error:type_name -> game_proto.ErrorId
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_error_proto_init() }
func file_error_proto_init() {
	if File_error_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_error_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgErrorAck); i {
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
			RawDescriptor: file_error_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_error_proto_goTypes,
		DependencyIndexes: file_error_proto_depIdxs,
		EnumInfos:         file_error_proto_enumTypes,
		MessageInfos:      file_error_proto_msgTypes,
	}.Build()
	File_error_proto = out.File
	file_error_proto_rawDesc = nil
	file_error_proto_goTypes = nil
	file_error_proto_depIdxs = nil
}
