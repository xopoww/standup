// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: standup_service.proto

package standup

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_standup_service_proto protoreflect.FileDescriptor

var file_standup_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x75, 0x70, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x78, 0x6f, 0x70, 0x6f, 0x77, 0x77, 0x2e,
	0x73, 0x74, 0x61, 0x6e, 0x64, 0x75, 0x70, 0x1a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x69, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x6e, 0x64, 0x75,
	0x70, 0x12, 0x5e, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x24, 0x2e, 0x78, 0x6f, 0x70, 0x6f, 0x77, 0x77, 0x2e, 0x73, 0x74, 0x61, 0x6e,
	0x64, 0x75, 0x70, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x78, 0x6f, 0x70, 0x6f, 0x77,
	0x77, 0x2e, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x75, 0x70, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x78, 0x6f, 0x70, 0x6f, 0x77, 0x77, 0x2f, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x75, 0x70, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x75, 0x70, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_standup_service_proto_goTypes = []interface{}{
	(*CreateMessageRequest)(nil),  // 0: xopoww.standup.CreateMessageRequest
	(*CreateMessageResponse)(nil), // 1: xopoww.standup.CreateMessageResponse
}
var file_standup_service_proto_depIdxs = []int32{
	0, // 0: xopoww.standup.Standup.CreateMessage:input_type -> xopoww.standup.CreateMessageRequest
	1, // 1: xopoww.standup.Standup.CreateMessage:output_type -> xopoww.standup.CreateMessageResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_standup_service_proto_init() }
func file_standup_service_proto_init() {
	if File_standup_service_proto != nil {
		return
	}
	file_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_standup_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_standup_service_proto_goTypes,
		DependencyIndexes: file_standup_service_proto_depIdxs,
	}.Build()
	File_standup_service_proto = out.File
	file_standup_service_proto_rawDesc = nil
	file_standup_service_proto_goTypes = nil
	file_standup_service_proto_depIdxs = nil
}
