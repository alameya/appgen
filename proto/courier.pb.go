// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.0
// source: proto/courier.proto

package proto

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

type CreateCourierRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateCourierRequest) Reset() {
	*x = CreateCourierRequest{}
	mi := &file_proto_courier_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCourierRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCourierRequest) ProtoMessage() {}

func (x *CreateCourierRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCourierRequest.ProtoReflect.Descriptor instead.
func (*CreateCourierRequest) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{0}
}

func (x *CreateCourierRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCourierRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *CreateCourierRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateCourierRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetCourierRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCourierRequest) Reset() {
	*x = GetCourierRequest{}
	mi := &file_proto_courier_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCourierRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCourierRequest) ProtoMessage() {}

func (x *GetCourierRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCourierRequest.ProtoReflect.Descriptor instead.
func (*GetCourierRequest) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{1}
}

func (x *GetCourierRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ListCourierRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCourierRequest) Reset() {
	*x = ListCourierRequest{}
	mi := &file_proto_courier_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCourierRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCourierRequest) ProtoMessage() {}

func (x *ListCourierRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCourierRequest.ProtoReflect.Descriptor instead.
func (*ListCourierRequest) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{2}
}

type ListCourierResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Courier             `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCourierResponse) Reset() {
	*x = ListCourierResponse{}
	mi := &file_proto_courier_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCourierResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCourierResponse) ProtoMessage() {}

func (x *ListCourierResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCourierResponse.ProtoReflect.Descriptor instead.
func (*ListCourierResponse) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{3}
}

func (x *ListCourierResponse) GetItems() []*Courier {
	if x != nil {
		return x.Items
	}
	return nil
}

type Courier struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	Email         string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Status        string                 `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Courier) Reset() {
	*x = Courier{}
	mi := &file_proto_courier_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Courier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Courier) ProtoMessage() {}

func (x *Courier) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Courier.ProtoReflect.Descriptor instead.
func (*Courier) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{4}
}

func (x *Courier) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Courier) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Courier) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Courier) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Courier) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type Location struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Latitude      float64                `protobuf:"fixed64,2,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude     float64                `protobuf:"fixed64,3,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Address       string                 `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Phone         string                 `protobuf:"bytes,6,opt,name=phone,proto3" json:"phone,omitempty"`
	Email         string                 `protobuf:"bytes,7,opt,name=email,proto3" json:"email,omitempty"`
	Status        string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Location) Reset() {
	*x = Location{}
	mi := &file_proto_courier_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_proto_courier_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_proto_courier_proto_rawDescGZIP(), []int{5}
}

func (x *Location) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Location) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Location) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *Location) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Location) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Location) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Location) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Location) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_courier_proto protoreflect.FileDescriptor

var file_proto_courier_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x22, 0x6e,
	0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x23,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3d, 0x0a, 0x13, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x26, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65,
	0x72, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x71, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x72,
	0x69, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xc6, 0x01, 0x0a, 0x08,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x32, 0xc9, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x1d, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x10, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x75, 0x72, 0x69,
	0x65, 0x72, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x63, 0x6f,
	0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x04, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x63, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x14, 0x5a, 0x12, 0x61, 0x70, 0x70, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_courier_proto_rawDescOnce sync.Once
	file_proto_courier_proto_rawDescData []byte
)

func file_proto_courier_proto_rawDescGZIP() []byte {
	file_proto_courier_proto_rawDescOnce.Do(func() {
		file_proto_courier_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_courier_proto_rawDesc), len(file_proto_courier_proto_rawDesc)))
	})
	return file_proto_courier_proto_rawDescData
}

var file_proto_courier_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_courier_proto_goTypes = []any{
	(*CreateCourierRequest)(nil), // 0: courier.CreateCourierRequest
	(*GetCourierRequest)(nil),    // 1: courier.GetCourierRequest
	(*ListCourierRequest)(nil),   // 2: courier.ListCourierRequest
	(*ListCourierResponse)(nil),  // 3: courier.ListCourierResponse
	(*Courier)(nil),              // 4: courier.Courier
	(*Location)(nil),             // 5: courier.Location
}
var file_proto_courier_proto_depIdxs = []int32{
	4, // 0: courier.ListCourierResponse.items:type_name -> courier.Courier
	0, // 1: courier.CourierService.Create:input_type -> courier.CreateCourierRequest
	1, // 2: courier.CourierService.Get:input_type -> courier.GetCourierRequest
	2, // 3: courier.CourierService.List:input_type -> courier.ListCourierRequest
	4, // 4: courier.CourierService.Create:output_type -> courier.Courier
	4, // 5: courier.CourierService.Get:output_type -> courier.Courier
	3, // 6: courier.CourierService.List:output_type -> courier.ListCourierResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_courier_proto_init() }
func file_proto_courier_proto_init() {
	if File_proto_courier_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_courier_proto_rawDesc), len(file_proto_courier_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_courier_proto_goTypes,
		DependencyIndexes: file_proto_courier_proto_depIdxs,
		MessageInfos:      file_proto_courier_proto_msgTypes,
	}.Build()
	File_proto_courier_proto = out.File
	file_proto_courier_proto_goTypes = nil
	file_proto_courier_proto_depIdxs = nil
}
