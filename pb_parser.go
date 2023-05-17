package protobuf

import (
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func ParsePBFileDesc(file []byte, ops ...OptionFunc) (*descriptor.FileDescriptorProto, error) {
	rr := new(descriptor.FileDescriptorProto)
	err := cgo_parse_pb(file, LoadOptions(append(ops, WithMessageType(MessageType_PB))...), func(desc []byte) error {
		return proto.Unmarshal(desc, rr)
	})
	if err != nil {
		return nil, err
	}
	return rr, nil
}

func ParsePBMultiFileDesc(idl *IDLConfig, ops ...OptionFunc) (*descriptor.FileDescriptorSet, error) {
	rr := new(descriptor.FileDescriptorSet)
	if err := cgo_parse_multi_pb(idl, LoadOptions(append(ops, WithMessageType(MessageType_PB))...), func(desc []byte) error {
		return proto.Unmarshal(desc, rr)
	}); err != nil {
		return nil, err
	}
	return rr, nil
}
