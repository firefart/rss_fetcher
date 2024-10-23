// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: rss.proto

package pb

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

type Rss struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastRun int64            `protobuf:"varint,1,opt,name=LastRun,proto3" json:"LastRun,omitempty"`
	Feeds   map[string]int64 `protobuf:"bytes,2,rep,name=Feeds,proto3" json:"Feeds,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Rss) Reset() {
	*x = Rss{}
	mi := &file_rss_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rss) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rss) ProtoMessage() {}

func (x *Rss) ProtoReflect() protoreflect.Message {
	mi := &file_rss_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rss.ProtoReflect.Descriptor instead.
func (*Rss) Descriptor() ([]byte, []int) {
	return file_rss_proto_rawDescGZIP(), []int{0}
}

func (x *Rss) GetLastRun() int64 {
	if x != nil {
		return x.LastRun
	}
	return 0
}

func (x *Rss) GetFeeds() map[string]int64 {
	if x != nil {
		return x.Feeds
	}
	return nil
}

var File_rss_proto protoreflect.FileDescriptor

var file_rss_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72, 0x73, 0x73,
	0x22, 0x84, 0x01, 0x0a, 0x03, 0x52, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x61, 0x73, 0x74,
	0x52, 0x75, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4c, 0x61, 0x73, 0x74, 0x52,
	0x75, 0x6e, 0x12, 0x29, 0x0a, 0x05, 0x46, 0x65, 0x65, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x72, 0x73, 0x73, 0x2e, 0x52, 0x73, 0x73, 0x2e, 0x46, 0x65, 0x65, 0x64,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x46, 0x65, 0x65, 0x64, 0x73, 0x1a, 0x38, 0x0a,
	0x0a, 0x46, 0x65, 0x65, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x5a, 0x0b, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rss_proto_rawDescOnce sync.Once
	file_rss_proto_rawDescData = file_rss_proto_rawDesc
)

func file_rss_proto_rawDescGZIP() []byte {
	file_rss_proto_rawDescOnce.Do(func() {
		file_rss_proto_rawDescData = protoimpl.X.CompressGZIP(file_rss_proto_rawDescData)
	})
	return file_rss_proto_rawDescData
}

var file_rss_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rss_proto_goTypes = []any{
	(*Rss)(nil), // 0: rss.Rss
	nil,         // 1: rss.Rss.FeedsEntry
}
var file_rss_proto_depIdxs = []int32{
	1, // 0: rss.Rss.Feeds:type_name -> rss.Rss.FeedsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rss_proto_init() }
func file_rss_proto_init() {
	if File_rss_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rss_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rss_proto_goTypes,
		DependencyIndexes: file_rss_proto_depIdxs,
		MessageInfos:      file_rss_proto_msgTypes,
	}.Build()
	File_rss_proto = out.File
	file_rss_proto_rawDesc = nil
	file_rss_proto_goTypes = nil
	file_rss_proto_depIdxs = nil
}
