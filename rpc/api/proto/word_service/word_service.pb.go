// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.15.8
// source: api/proto/word_service/word_service.proto

package word_service

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

type DefineWordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lexicon     string   `protobuf:"bytes,1,opt,name=lexicon,proto3" json:"lexicon,omitempty"`
	Words       []string `protobuf:"bytes,2,rep,name=words,proto3" json:"words,omitempty"`
	Definitions bool     `protobuf:"varint,3,opt,name=definitions,proto3" json:"definitions,omitempty"` // pass true to retrieve definitions
}

func (x *DefineWordsRequest) Reset() {
	*x = DefineWordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_word_service_word_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefineWordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefineWordsRequest) ProtoMessage() {}

func (x *DefineWordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_word_service_word_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefineWordsRequest.ProtoReflect.Descriptor instead.
func (*DefineWordsRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_word_service_word_service_proto_rawDescGZIP(), []int{0}
}

func (x *DefineWordsRequest) GetLexicon() string {
	if x != nil {
		return x.Lexicon
	}
	return ""
}

func (x *DefineWordsRequest) GetWords() []string {
	if x != nil {
		return x.Words
	}
	return nil
}

func (x *DefineWordsRequest) GetDefinitions() bool {
	if x != nil {
		return x.Definitions
	}
	return false
}

type DefineWordsResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	D string `protobuf:"bytes,1,opt,name=d,proto3" json:"d,omitempty"`  // definitions, not "" iff (valid and requesting definitions)
	V bool   `protobuf:"varint,2,opt,name=v,proto3" json:"v,omitempty"` // true iff valid
}

func (x *DefineWordsResult) Reset() {
	*x = DefineWordsResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_word_service_word_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefineWordsResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefineWordsResult) ProtoMessage() {}

func (x *DefineWordsResult) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_word_service_word_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefineWordsResult.ProtoReflect.Descriptor instead.
func (*DefineWordsResult) Descriptor() ([]byte, []int) {
	return file_api_proto_word_service_word_service_proto_rawDescGZIP(), []int{1}
}

func (x *DefineWordsResult) GetD() string {
	if x != nil {
		return x.D
	}
	return ""
}

func (x *DefineWordsResult) GetV() bool {
	if x != nil {
		return x.V
	}
	return false
}

type DefineWordsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results map[string]*DefineWordsResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DefineWordsResponse) Reset() {
	*x = DefineWordsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_word_service_word_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefineWordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefineWordsResponse) ProtoMessage() {}

func (x *DefineWordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_word_service_word_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefineWordsResponse.ProtoReflect.Descriptor instead.
func (*DefineWordsResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_word_service_word_service_proto_rawDescGZIP(), []int{2}
}

func (x *DefineWordsResponse) GetResults() map[string]*DefineWordsResult {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_api_proto_word_service_word_service_proto protoreflect.FileDescriptor

var file_api_proto_word_service_word_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x64,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x6f, 0x72,
	0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x66, 0x0a, 0x12, 0x44, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x6c, 0x65, 0x78, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x65, 0x78, 0x69, 0x63, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6f, 0x72,
	0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x2f, 0x0a, 0x11, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x01, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x01, 0x76, 0x22, 0xbc, 0x01, 0x0a, 0x13, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72,
	0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x77, 0x6f,
	0x72, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x66, 0x69, 0x6e,
	0x65, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x1a, 0x5b, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x35, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x32, 0x61, 0x0a, 0x0b, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x52, 0x0a, 0x0b, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x12,
	0x20, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44,
	0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x6d, 0x69, 0x6e, 0x6f, 0x31, 0x34, 0x2f, 0x6c, 0x69, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_word_service_word_service_proto_rawDescOnce sync.Once
	file_api_proto_word_service_word_service_proto_rawDescData = file_api_proto_word_service_word_service_proto_rawDesc
)

func file_api_proto_word_service_word_service_proto_rawDescGZIP() []byte {
	file_api_proto_word_service_word_service_proto_rawDescOnce.Do(func() {
		file_api_proto_word_service_word_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_word_service_word_service_proto_rawDescData)
	})
	return file_api_proto_word_service_word_service_proto_rawDescData
}

var file_api_proto_word_service_word_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proto_word_service_word_service_proto_goTypes = []interface{}{
	(*DefineWordsRequest)(nil),  // 0: word_service.DefineWordsRequest
	(*DefineWordsResult)(nil),   // 1: word_service.DefineWordsResult
	(*DefineWordsResponse)(nil), // 2: word_service.DefineWordsResponse
	nil,                         // 3: word_service.DefineWordsResponse.ResultsEntry
}
var file_api_proto_word_service_word_service_proto_depIdxs = []int32{
	3, // 0: word_service.DefineWordsResponse.results:type_name -> word_service.DefineWordsResponse.ResultsEntry
	1, // 1: word_service.DefineWordsResponse.ResultsEntry.value:type_name -> word_service.DefineWordsResult
	0, // 2: word_service.WordService.DefineWords:input_type -> word_service.DefineWordsRequest
	2, // 3: word_service.WordService.DefineWords:output_type -> word_service.DefineWordsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_word_service_word_service_proto_init() }
func file_api_proto_word_service_word_service_proto_init() {
	if File_api_proto_word_service_word_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_word_service_word_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefineWordsRequest); i {
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
		file_api_proto_word_service_word_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefineWordsResult); i {
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
		file_api_proto_word_service_word_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefineWordsResponse); i {
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
			RawDescriptor: file_api_proto_word_service_word_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_word_service_word_service_proto_goTypes,
		DependencyIndexes: file_api_proto_word_service_word_service_proto_depIdxs,
		MessageInfos:      file_api_proto_word_service_word_service_proto_msgTypes,
	}.Build()
	File_api_proto_word_service_word_service_proto = out.File
	file_api_proto_word_service_word_service_proto_rawDesc = nil
	file_api_proto_word_service_word_service_proto_goTypes = nil
	file_api_proto_word_service_word_service_proto_depIdxs = nil
}
