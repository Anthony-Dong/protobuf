package benchmark

import (
	"errors"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/anthony-dong/protobuf"
	"github.com/jhump/protoreflect/desc/protoparse"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// Benchmark_ParsePBFileDesc_Cgo
// use "github.com/anthony-dong/protobuf"
func Benchmark_ParsePBFileDesc_Cgo(b *testing.B) {
	b.StopTimer()
	mainIdl := loadIdl(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rr, err := protobuf.ParsePBFileDesc(mainIdl, protobuf.WithRequireSyntaxIdentifier())
		if err != nil {
			b.Fatal(err)
		}
		if rr == nil {
			b.Fatal("must error")
		}
	}
}

// Benchmark_ParsePBFileDesc_Jhump
// use "github.com/jhump/protoreflect"
func Benchmark_ParsePBFileDesc_Jhump(b *testing.B) {
	b.StopTimer()
	mainIdl := string(loadIdl(b))
	mainIdlName := "main.proto"
	idlMap := map[string]string{mainIdlName: mainIdl}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rr, err := ParsePBFileDescJhump(idlMap, mainIdlName)
		if err != nil {
			b.Fatal(err)
		}
		if rr == nil {
			b.Fatal("must error")
		}
	}
}

func ParsePBFileDescJhump(idls map[string]string, main string) (*descriptor.FileDescriptorProto, error) {
	var pbParser protoparse.Parser
	pbParser.Accessor = protoparse.FileContentsFromMap(idls)
	//pbParser.IncludeSourceCodeInfo = true // 关闭这个可以降低内存开销. 尤其对于大型idl来说!
	pbParser.ValidateUnlinkedFiles = true
	pbParser.InterpretOptionsInUnlinkedFiles = true
	fds, err := pbParser.ParseFiles(main)
	if err != nil {
		return nil, err
	}
	for _, fd := range fds {
		return fd.AsFileDescriptorProto(), nil
	}
	return nil, errors.New("no file desc")
}

var initOnce sync.Once
var idl []byte

func loadIdl(t testing.TB) []byte {
	initOnce.Do(func() {
		file, err := ioutil.ReadFile("../test/api.proto")
		if err != nil {
			t.Fatal(err)
		}
		idl = file
	})
	return idl
}
