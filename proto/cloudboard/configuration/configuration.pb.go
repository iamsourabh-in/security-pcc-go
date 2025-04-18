// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: proto/cloudboard/configuration/configuration.proto

package configuration

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RegisterRequest is sent by a client to register for updates.
type RegisterRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Serialized registration payload (e.g. domain, client info).
	Registration  []byte `protobuf:"bytes,1,opt,name=registration,proto3" json:"registration,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetRegistration() []byte {
	if x != nil {
		return x.Registration
	}
	return nil
}

// ConfigurationUpdate contains the update data to apply.
type ConfigurationUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UpdatePayload []byte                 `protobuf:"bytes,1,opt,name=update_payload,json=updatePayload,proto3" json:"update_payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConfigurationUpdate) Reset() {
	*x = ConfigurationUpdate{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigurationUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationUpdate) ProtoMessage() {}

func (x *ConfigurationUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationUpdate.ProtoReflect.Descriptor instead.
func (*ConfigurationUpdate) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{1}
}

func (x *ConfigurationUpdate) GetUpdatePayload() []byte {
	if x != nil {
		return x.UpdatePayload
	}
	return nil
}

// ApplySuccessRequest acknowledges a successful apply.
type ApplySuccessRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	SuccessPayload []byte                 `protobuf:"bytes,1,opt,name=success_payload,json=successPayload,proto3" json:"success_payload,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ApplySuccessRequest) Reset() {
	*x = ApplySuccessRequest{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ApplySuccessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplySuccessRequest) ProtoMessage() {}

func (x *ApplySuccessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplySuccessRequest.ProtoReflect.Descriptor instead.
func (*ApplySuccessRequest) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{2}
}

func (x *ApplySuccessRequest) GetSuccessPayload() []byte {
	if x != nil {
		return x.SuccessPayload
	}
	return nil
}

// ApplyFailureRequest acknowledges a failed apply.
type ApplyFailureRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	FailurePayload []byte                 `protobuf:"bytes,1,opt,name=failure_payload,json=failurePayload,proto3" json:"failure_payload,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ApplyFailureRequest) Reset() {
	*x = ApplyFailureRequest{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ApplyFailureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyFailureRequest) ProtoMessage() {}

func (x *ApplyFailureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyFailureRequest.ProtoReflect.Descriptor instead.
func (*ApplyFailureRequest) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{3}
}

func (x *ApplyFailureRequest) GetFailurePayload() []byte {
	if x != nil {
		return x.FailurePayload
	}
	return nil
}

// VersionInfoResponse returns info about the current config version.
type VersionInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VersionInfo   []byte                 `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersionInfoResponse) Reset() {
	*x = VersionInfoResponse{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersionInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionInfoResponse) ProtoMessage() {}

func (x *VersionInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionInfoResponse.ProtoReflect.Descriptor instead.
func (*VersionInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{4}
}

func (x *VersionInfoResponse) GetVersionInfo() []byte {
	if x != nil {
		return x.VersionInfo
	}
	return nil
}

// EmptyRequest is a placeholder for empty requests.
type EmptyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{5}
}

// EmptyResponse is a placeholder for empty responses.
type EmptyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmptyResponse) Reset() {
	*x = EmptyResponse{}
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmptyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyResponse) ProtoMessage() {}

func (x *EmptyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cloudboard_configuration_configuration_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyResponse.ProtoReflect.Descriptor instead.
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP(), []int{6}
}

var File_proto_cloudboard_configuration_configuration_proto protoreflect.FileDescriptor

const file_proto_cloudboard_configuration_configuration_proto_rawDesc = "" +
	"\n" +
	"2proto/cloudboard/configuration/configuration.proto\x12\x18cloudboard.configuration\"5\n" +
	"\x0fRegisterRequest\x12\"\n" +
	"\fregistration\x18\x01 \x01(\fR\fregistration\"<\n" +
	"\x13ConfigurationUpdate\x12%\n" +
	"\x0eupdate_payload\x18\x01 \x01(\fR\rupdatePayload\">\n" +
	"\x13ApplySuccessRequest\x12'\n" +
	"\x0fsuccess_payload\x18\x01 \x01(\fR\x0esuccessPayload\">\n" +
	"\x13ApplyFailureRequest\x12'\n" +
	"\x0ffailure_payload\x18\x01 \x01(\fR\x0efailurePayload\"8\n" +
	"\x13VersionInfoResponse\x12!\n" +
	"\fversion_info\x18\x01 \x01(\fR\vversionInfo\"\x0e\n" +
	"\fEmptyRequest\"\x0f\n" +
	"\rEmptyResponse2\xe1\x03\n" +
	"\rConfiguration\x12d\n" +
	"\bRegister\x12).cloudboard.configuration.RegisterRequest\x1a-.cloudboard.configuration.ConfigurationUpdate\x12z\n" +
	" SuccessfullyAppliedConfiguration\x12-.cloudboard.configuration.ApplySuccessRequest\x1a'.cloudboard.configuration.EmptyResponse\x12t\n" +
	"\x1aFailedToApplyConfiguration\x12-.cloudboard.configuration.ApplyFailureRequest\x1a'.cloudboard.configuration.EmptyResponse\x12x\n" +
	"\x1fCurrentConfigurationVersionInfo\x12&.cloudboard.configuration.EmptyRequest\x1a-.cloudboard.configuration.VersionInfoResponseBIZGgithub.com/iamsourabh-in/security-pcc-go/proto/cloudboard/configurationb\x06proto3"

var (
	file_proto_cloudboard_configuration_configuration_proto_rawDescOnce sync.Once
	file_proto_cloudboard_configuration_configuration_proto_rawDescData []byte
)

func file_proto_cloudboard_configuration_configuration_proto_rawDescGZIP() []byte {
	file_proto_cloudboard_configuration_configuration_proto_rawDescOnce.Do(func() {
		file_proto_cloudboard_configuration_configuration_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_cloudboard_configuration_configuration_proto_rawDesc), len(file_proto_cloudboard_configuration_configuration_proto_rawDesc)))
	})
	return file_proto_cloudboard_configuration_configuration_proto_rawDescData
}

var file_proto_cloudboard_configuration_configuration_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_cloudboard_configuration_configuration_proto_goTypes = []any{
	(*RegisterRequest)(nil),     // 0: cloudboard.configuration.RegisterRequest
	(*ConfigurationUpdate)(nil), // 1: cloudboard.configuration.ConfigurationUpdate
	(*ApplySuccessRequest)(nil), // 2: cloudboard.configuration.ApplySuccessRequest
	(*ApplyFailureRequest)(nil), // 3: cloudboard.configuration.ApplyFailureRequest
	(*VersionInfoResponse)(nil), // 4: cloudboard.configuration.VersionInfoResponse
	(*EmptyRequest)(nil),        // 5: cloudboard.configuration.EmptyRequest
	(*EmptyResponse)(nil),       // 6: cloudboard.configuration.EmptyResponse
}
var file_proto_cloudboard_configuration_configuration_proto_depIdxs = []int32{
	0, // 0: cloudboard.configuration.Configuration.Register:input_type -> cloudboard.configuration.RegisterRequest
	2, // 1: cloudboard.configuration.Configuration.SuccessfullyAppliedConfiguration:input_type -> cloudboard.configuration.ApplySuccessRequest
	3, // 2: cloudboard.configuration.Configuration.FailedToApplyConfiguration:input_type -> cloudboard.configuration.ApplyFailureRequest
	5, // 3: cloudboard.configuration.Configuration.CurrentConfigurationVersionInfo:input_type -> cloudboard.configuration.EmptyRequest
	1, // 4: cloudboard.configuration.Configuration.Register:output_type -> cloudboard.configuration.ConfigurationUpdate
	6, // 5: cloudboard.configuration.Configuration.SuccessfullyAppliedConfiguration:output_type -> cloudboard.configuration.EmptyResponse
	6, // 6: cloudboard.configuration.Configuration.FailedToApplyConfiguration:output_type -> cloudboard.configuration.EmptyResponse
	4, // 7: cloudboard.configuration.Configuration.CurrentConfigurationVersionInfo:output_type -> cloudboard.configuration.VersionInfoResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_cloudboard_configuration_configuration_proto_init() }
func file_proto_cloudboard_configuration_configuration_proto_init() {
	if File_proto_cloudboard_configuration_configuration_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_cloudboard_configuration_configuration_proto_rawDesc), len(file_proto_cloudboard_configuration_configuration_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_cloudboard_configuration_configuration_proto_goTypes,
		DependencyIndexes: file_proto_cloudboard_configuration_configuration_proto_depIdxs,
		MessageInfos:      file_proto_cloudboard_configuration_configuration_proto_msgTypes,
	}.Build()
	File_proto_cloudboard_configuration_configuration_proto = out.File
	file_proto_cloudboard_configuration_configuration_proto_goTypes = nil
	file_proto_cloudboard_configuration_configuration_proto_depIdxs = nil
}
