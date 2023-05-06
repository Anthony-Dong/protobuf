package protobuf

import (
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func ParsePBFileDesc(file []byte, ops ...OptionFunc) (*descriptor.FileDescriptorProto, error) {
	pbBinary, err := cgo_parse_pb(file, append(ops, WithMessageType(MessageType_PB))...)
	if err != nil {
		return nil, err
	}
	rr := new(descriptor.FileDescriptorProto)
	if err := proto.Unmarshal(pbBinary, rr); err != nil {
		return nil, err
	}
	return rr, nil
}
