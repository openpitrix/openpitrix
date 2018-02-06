// Code generated by protoc-gen-go. DO NOT EDIT.
// source: annotations.proto

/*
Package openpitrix is a generated protocol buffer package.

It is generated from these files:
	annotations.proto
	app.proto
	app_runtime.proto
	cluster.proto
	metadata.proto
	repo.proto

It has these top-level messages:
	App
	AppId
	AppListRequest
	AppListResponse
	AppRuntime
	AppRuntimeLabel
	AppRuntimeId
	AppRuntimeListRequest
	AppRuntimeListResponse
	AppRuntimePluginInfo
	AppRuntimePluginInput
	AppRuntimePluginOutput
	Cluster
	Clusters
	ClusterNode
	ClusterNodes
	ClusterId
	ClusterIds
	ClusterListRequest
	ClusterListResponse
	ClusterNodeId
	ClusterNodeIds
	ClusterNodeListRequest
	ClusterNodeListResponse
	Const
	Default
	DBAppTableSchema
	DBRuntimeTableSchema
	DBClusterTableSchema
	DBClusterNodeTableSchema
	DBRepoTableSchema
	Repo
	RepoLabel
	RepoSelector
	RepoId
	RepoListRequest
	RepoListResponse
*/
package openpitrix

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import grpc_gateway_protoc_gen_swagger_options "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import google_protobuf1 "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

var E_Openapiv2FieldSchema = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf1.FieldOptions)(nil),
	ExtensionType: (*grpc_gateway_protoc_gen_swagger_options.JSONSchema)(nil),
	Field:         1042,
	Name:          "openpitrix.openapiv2_field_schema",
	Tag:           "bytes,1042,opt,name=openapiv2_field_schema,json=openapiv2FieldSchema",
	Filename:      "annotations.proto",
}

func init() {
	proto.RegisterExtension(E_Openapiv2FieldSchema)
}

func init() { proto.RegisterFile("annotations.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x8e, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0xe9, 0x4d, 0xea, 0xc9, 0x22, 0x22, 0x05, 0xa1, 0x47, 0x11, 0x3a, 0x81, 0x7a, 0xf3,
	0x01, 0x3c, 0x78, 0xb0, 0x60, 0x1f, 0x20, 0xc4, 0x74, 0x3a, 0x06, 0x6a, 0x26, 0x24, 0xd1, 0xee,
	0x3e, 0xc2, 0x5e, 0xf7, 0x89, 0x97, 0x26, 0xdd, 0x1e, 0x33, 0xf9, 0xbe, 0xff, 0xff, 0xcb, 0x3b,
	0x65, 0x2d, 0x47, 0x15, 0x0d, 0xdb, 0x00, 0xce, 0x73, 0xe4, 0xaa, 0x64, 0x87, 0xd6, 0x99, 0xe8,
	0xcd, 0xa1, 0x7e, 0x49, 0x27, 0xdd, 0x12, 0xda, 0x36, 0x2c, 0x8a, 0x08, 0xbd, 0x60, 0x97, 0x68,
	0xb1, 0x62, 0xca, 0x99, 0xff, 0x2e, 0x7b, 0x75, 0x43, 0xcc, 0x34, 0xa3, 0x48, 0xaf, 0xef, 0xbf,
	0x49, 0x8c, 0x18, 0xb4, 0x37, 0x2e, 0xb2, 0xcf, 0xc4, 0xdb, 0xa9, 0x28, 0x1f, 0x76, 0x4b, 0x4e,
	0x06, 0xe7, 0x51, 0x06, 0xfd, 0x83, 0xbf, 0xaa, 0x7a, 0x82, 0x6c, 0xc3, 0xd5, 0x86, 0xf7, 0xf5,
	0xbb, 0xcf, 0x5d, 0x8f, 0xe7, 0x9b, 0xa6, 0x78, 0xbe, 0xed, 0x5e, 0x81, 0xbc, 0xd3, 0x40, 0x2a,
	0xe2, 0xa2, 0x8e, 0x99, 0xd5, 0x92, 0xd0, 0xca, 0x6d, 0x1c, 0x6c, 0xe3, 0xe0, 0x63, 0xe8, 0x3f,
	0x87, 0x14, 0xfd, 0x75, 0xbf, 0x57, 0xa6, 0xc8, 0x7c, 0xbd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x2d,
	0xca, 0xff, 0x5f, 0xf9, 0x00, 0x00, 0x00,
}
