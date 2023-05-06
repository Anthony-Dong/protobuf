package protobuf

import (
	"reflect"
	"unsafe"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnsafeBytes(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func MessageToJson(v proto.Message, pretty ...bool) string {
	ops := protojson.MarshalOptions{Multiline: len(pretty) > 0 && pretty[0]}
	return ops.Format(v)
}
