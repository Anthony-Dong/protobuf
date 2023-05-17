package protobuf

import (
	"io/ioutil"
	"sync"
	"testing"
)

func Test_cgo_parse_pb(t *testing.T) {
	test_parse_pb(t)
	test_parse_pb_json(t)
}

var testJsonOps = &Option{
	MessageType:             MessageType_JSON,
	RequireSyntaxIdentifier: true,
}

var testPbOps = &Option{
	MessageType:             MessageType_PB,
	RequireSyntaxIdentifier: true,
}

func test_parse_pb(t testing.TB) {
	if err := cgo_parse_pb(loadIdl(t), testPbOps, func(desc []byte) error {
		if len(desc) == 0 {
			t.Fatal("error..")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func test_parse_pb_json(t testing.TB) {
	err := cgo_parse_pb(loadIdl(t), testJsonOps, func(desc []byte) error {
		if len(desc) == 0 {
			t.Fatal("error..")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Benchmark_cgo_parse_pb_pb(b *testing.B) {
	b.StopTimer()
	_ = loadIdl(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		test_parse_pb(b)
	}
}

func Benchmark_cgo_parse_pb_json(b *testing.B) {
	b.StopTimer()
	_ = loadIdl(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		test_parse_pb_json(b)
	}
}

func Test_cgo_parse_pb_error(t *testing.T) {
	var pb []byte
	err := cgo_parse_pb([]byte(`hello world`), testJsonOps, func(desc []byte) error {
		return nil
	})
	if err != nil {
		messages := err.(ErrorMessages)
		for _, elem := range messages {
			t.Logf("%#v\n", elem)
		}
		t.Log(err.Error())
	}
	if err == nil {
		t.Fatal("must error")
	}
	if pb != nil {
		t.Fatal("pb error")
	}
}

var initOnce sync.Once
var idl []byte

func loadIdl(t testing.TB) []byte {
	initOnce.Do(func() {
		file, err := ioutil.ReadFile("internal/idl/test.proto")
		if err != nil {
			t.Fatal(err)
		}
		idl = file
	})
	return idl
}

var idlConfig *IDLConfig
var idlConfigInit sync.Once

func loadIdls(t testing.TB) *IDLConfig {
	idlConfigInit.Do(func() {
		tree, err := NewProtobufDiskSourceTree("internal/idl")
		if err != nil {
			t.Fatal(err)
		}
		idlConfig = new(IDLConfig)
		idlConfig.IDLs = tree
		idlConfig.Main = "service/im.proto"
		idlConfig.IncludePath = []string{"desc", "."}
	})
	return idlConfig
}
