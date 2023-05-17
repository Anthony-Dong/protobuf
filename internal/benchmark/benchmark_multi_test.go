package benchmark

import (
	"errors"
	"sync"
	"testing"

	"github.com/jhump/protoreflect/desc"

	"github.com/anthony-dong/protobuf"
	"github.com/jhump/protoreflect/desc/protoparse"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

// Benchmark_ParsePBMultiFileDesc_Cgo
// use "github.com/anthony-dong/protobuf"
func Benchmark_ParsePBMultiFileDesc_Cgo(b *testing.B) {
	b.StopTimer()
	mainIdl := loadIdls(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rr, err := protobuf.ParsePBMultiFileDesc(mainIdl, protobuf.WithRequireSyntaxIdentifier(), protobuf.WithJsonTag(), protobuf.WithGoogleProtobuf())
		if err != nil {
			b.Fatal(err)
		}
		if rr == nil {
			b.Fatal("must error")
		}
	}
}

// Benchmark_ParsePBMultiFileDesc_Jhump
// use "github.com/jhump/protoreflect"
func Benchmark_ParsePBMultiFileDesc_Jhump(b *testing.B) {
	b.StopTimer()
	mainIdl := loadIdls(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rr, err := ParsePBMultiFileDesc(mainIdl)
		if err != nil {
			b.Fatal(err)
		}
		if rr == nil {
			b.Fatal("must error")
		}
	}
}

func TestParsePBMultiFileDesc(t *testing.T) {
	t.Run("jhump", func(t *testing.T) {
		result, err := ParsePBMultiFileDesc(loadIdls(t))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(protobuf.MessageToJson(result, true))
	})

	t.Run("cgo", func(t *testing.T) {
		result, err := protobuf.ParsePBMultiFileDesc(loadIdls(t), protobuf.WithRequireSyntaxIdentifier(), protobuf.WithJsonTag(), protobuf.WithGoogleProtobuf())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(protobuf.MessageToJson(result, true))
	})
}

func ParsePBMultiFileDesc(idls *protobuf.IDLConfig) (*descriptor.FileDescriptorSet, error) {
	var pbParser protoparse.Parser
	idlConfigMap := make(map[string]string, len(idlConfig.IDLs))
	for k, v := range idlConfig.IDLs {
		idlConfigMap[k] = protobuf.UnsafeString(v)
	}
	pbParser.Accessor = protoparse.FileContentsFromMap(idlConfigMap)
	//pbParser.IncludeSourceCodeInfo = true // 关闭这个可以降低内存开销. 尤其对于大型idl来说!
	pbParser.ValidateUnlinkedFiles = true
	pbParser.InterpretOptionsInUnlinkedFiles = true
	pbParser.ImportPaths = idls.IncludePath
	fds, err := pbParser.ParseFiles(idls.Main)
	if err != nil {
		return nil, err
	}
	result := descriptor.FileDescriptorSet{File: make([]*descriptor.FileDescriptorProto, 0, len(fds))}
	walk := make(map[string]bool, len(idls.IDLs))
	for _, fd := range fds {
		loadFds(&result, fd, walk)
	}
	if len(result.File) == 0 {
		return nil, errors.New("no file desc")
	}
	return &result, nil
}

func loadFds(set *descriptor.FileDescriptorSet, main *desc.FileDescriptor, walk map[string]bool) {
	if walk[main.GetName()] {
		return
	}
	walk[main.GetName()] = true

	set.File = append(set.File, main.AsFileDescriptorProto())

	for _, elem := range main.GetDependencies() {
		loadFds(set, elem, walk)
	}
}

var idlConfig *protobuf.IDLConfig
var idlConfigInit sync.Once

func loadIdls(t testing.TB) *protobuf.IDLConfig {
	idlConfigInit.Do(func() {
		tree, err := protobuf.NewProtobufDiskSourceTree("../idl")
		if err != nil {
			t.Fatal(err)
		}
		idlConfig = new(protobuf.IDLConfig)
		idlConfig.IDLs = tree
		idlConfig.Main = "service/im.proto"
		idlConfig.IncludePath = []string{"desc", "."}
	})
	return idlConfig
}
