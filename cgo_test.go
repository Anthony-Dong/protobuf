package protobuf

import (
	"io/ioutil"
	"sync"
	"testing"
)

func Test_cgo_parse_pb(t *testing.T) {
	test_parse_pb(t, MessageType_JSON)
	test_parse_pb(t, MessageType_PB)
}

func Test_cgo_parse_pb_error(t *testing.T) {
	pb, err := cgo_parse_pb([]byte(`hello world`), func(option *Option) {
		option.MessageType = MessageType_JSON
	})
	if err != nil {
		messages := err.(ErrorMessages)
		for _, elem := range messages {
			t.Logf("%#v\n", elem)
		}
	}
	if err == nil {
		t.Fatal("must error")
	}
	if pb != nil {
		t.Fatal("pb error")
	}
}

func test_parse_pb(t testing.TB, messageType MessageType) {
	cc, err := cgo_parse_pb(loadIdl(t), func(option *Option) {
		option.MessageType = messageType
		option.RequireSyntaxIdentifier = true
	})
	if err != nil {
		t.Fatal(err)
	}
	if cc == nil {
		t.Fatal("pb binary is nil")
	}
}

func Benchmark_cgo_parse_pb_pb(b *testing.B) {
	b.StopTimer()
	_ = loadIdl(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		test_parse_pb(b, MessageType_PB)
	}
}

func Benchmark_cgo_parse_pb_json(b *testing.B) {
	b.StopTimer()
	_ = loadIdl(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		test_parse_pb(b, MessageType_JSON)
	}
}

var initOnce sync.Once
var idl []byte

func loadIdl(t testing.TB) []byte {
	initOnce.Do(func() {
		file, err := ioutil.ReadFile("internal/test/api.proto")
		if err != nil {
			t.Fatal(err)
		}
		idl = file
	})
	return idl
}
