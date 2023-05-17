package main

import (
	"fmt"

	"github.com/anthony-dong/protobuf"
)

func main() {
	tree, err := protobuf.NewProtobufDiskSourceTree("internal/idl")
	if err != nil {
		panic(err)
	}
	idlConfig := new(protobuf.IDLConfig)
	idlConfig.IDLs = tree
	idlConfig.Main = "service/im.proto"
	idlConfig.IncludePath = []string{"desc", "."}

	desc, err := protobuf.ParsePBMultiFileDesc(idlConfig,
		protobuf.WithJsonTag(),
		protobuf.WithSourceCodeInfo(),
		protobuf.WithGoogleProtobuf(),
		protobuf.WithRequireSyntaxIdentifier(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(protobuf.MessageToJson(desc, true))
}
