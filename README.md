# Protobuf Parser

## 能力

1. 基于libprotoc 解析 proto 文件，解决一些开源的Go的Protobuf解析库和官方的库不一致的问题！

2. 支持的环境
- Darwin amd64
- Linux amd64

> 有其他环境需求可自行提交 MR，只要本地跑通测试即可！


## 使用

1. 下载依赖

```shell
CGO_ENABLED=1 go get -v github.com/anthony-dong/protobuf
```

2. 如何使用: [main.go](internal/example/main.go)

```go
package main

import (
	"log"

	"github.com/anthony-dong/protobuf"
)

func main() {
	file := []byte(`
syntax = "proto2";
package idl.model;
message Person {
  optional string name = 1;
  optional int32 id = 2;
  optional string email = 3;
  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
  }
  message PhoneNumber {
    optional string number = 1;
    optional PhoneType type = 2 [default = HOME];
  }
  repeated PhoneNumber phones = 4;
  map<string, Person> map_person = 5;
  optional bool status = 6;
}
`)
	desc, err := protobuf.ParsePBFileDesc(file, protobuf.WithRequireSyntaxIdentifier())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(protobuf.MessageToJson(desc))
}

// 运行: CGO_ENABLED=1 go run main.go
```

## 性能

测试文件: [api.proto](internal/test/api.proto)

```shell
➜  protobuf  make benchmark
go test -v -run=none -bench=Benchmark -memprofile mem.out  -benchmem  -count=5 .
goos: linux
goarch: amd64
pkg: github.com/anthony-dong/protobuf
cpu: Intel(R) Xeon(R) Platinum 8260 CPU @ 2.40GHz
Benchmark_cgo_parse_pb_pb
Benchmark_cgo_parse_pb_pb-8     	   25118	     49275 ns/op	    1432 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   25916	     47035 ns/op	    1432 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   25940	     46782 ns/op	    1432 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   25016	     48112 ns/op	    1432 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   25872	     46880 ns/op	    1432 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_json
Benchmark_cgo_parse_pb_json-8   	    6255	    192098 ns/op	    4120 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6019	    195550 ns/op	    4120 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6166	    194726 ns/op	    4120 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6153	    187322 ns/op	    4120 B/op	       3 allocs/op
Benchmark_cgo_parse_pb_json-8   	    5694	    187708 ns/op	    4120 B/op	       3 allocs/op
Benchmark_ParsePBFileDesc
Benchmark_ParsePBFileDesc-8     	   12237	     96951 ns/op	   22120 B/op	     549 allocs/op
Benchmark_ParsePBFileDesc-8     	   12291	     98412 ns/op	   22120 B/op	     549 allocs/op
Benchmark_ParsePBFileDesc-8     	   12217	     99688 ns/op	   22120 B/op	     549 allocs/op
Benchmark_ParsePBFileDesc-8     	   12262	     97279 ns/op	   22120 B/op	     549 allocs/op
Benchmark_ParsePBFileDesc-8     	   10000	    102308 ns/op	   22120 B/op	     549 allocs/op
PASS
ok  	github.com/anthony-dong/protobuf	24.224s
```