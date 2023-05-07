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

压测文件: [api.proto](internal/test/api.proto)

1. 解析PB性能 [cgo_test.go](./cgo_test.go)

```shell
go test -v -run=none -bench=Benchmark -memprofile mem.out  -benchmem  -count=5 .
goos: linux
goarch: amd64
pkg: github.com/anthony-dong/protobuf
cpu: Intel(R) Xeon(R) Platinum 8260 CPU @ 2.40GHz
Benchmark_cgo_parse_pb_pb
Benchmark_cgo_parse_pb_pb-8     	   25767	     42901 ns/op	    1424 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   27571	     41677 ns/op	    1424 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   28008	     43252 ns/op	    1424 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   28279	     42257 ns/op	    1424 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   27609	     46022 ns/op	    1424 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_json
Benchmark_cgo_parse_pb_json-8   	    5910	    178439 ns/op	    4112 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6274	    177877 ns/op	    4112 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6618	    175789 ns/op	    4112 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_json-8   	    5842	    185348 ns/op	    4112 B/op	       2 allocs/op
Benchmark_cgo_parse_pb_json-8   	    6403	    174213 ns/op	    4112 B/op	       2 allocs/op
PASS
ok  	github.com/anthony-dong/protobuf	13.813s
```

2. 对比 [github.com/jhump/protoreflect@v1.8.2](https://github.com/jhump/protoreflect/tree/v1.8.2 ) 解析库  [benchmark_test.go](internal/benchmark/benchmark_test.go)

> 原因是我司用的是 v1.8.2 版本, 高版本兼容性检测会高一些!

```shell
go test -v -run=none -bench=Benchmark -memprofile mem.out  -benchmem  -count=5 .
goos: linux
goarch: amd64
pkg: github.com/anthony-dong/protobuf/internal/benchmark
cpu: Intel(R) Xeon(R) Platinum 8260 CPU @ 2.40GHz
Benchmark_ParsePBFileDesc_Cgo
Benchmark_ParsePBFileDesc_Cgo-8     	   12895	     97235 ns/op	   22120 B/op	     548 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8     	   12114	     94504 ns/op	   22120 B/op	     548 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8     	   12573	     98407 ns/op	   22120 B/op	     548 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8     	   12562	     97684 ns/op	   22120 B/op	     548 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8     	   11883	     97633 ns/op	   22120 B/op	     548 allocs/op
Benchmark_ParsePBFileDesc_Jhump
Benchmark_ParsePBFileDesc_Jhump-8   	    4986	    243932 ns/op	  102432 B/op	    1580 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8   	    4754	    238628 ns/op	  102428 B/op	    1580 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8   	    5204	    240454 ns/op	  102424 B/op	    1580 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8   	    5052	    239961 ns/op	  102423 B/op	    1580 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8   	    4473	    247136 ns/op	  102431 B/op	    1580 allocs/op
PASS
ok  	github.com/anthony-dong/protobuf/internal/benchmark	16.892s
```

备注: 差异的主要原因在于 C++ 的内存分配性能要优于Go，对于Parser这种内存开销较大的业务逻辑，其次官方的 [protobuf](https://github.com/protocolbuffers/protobuf/tree/v3.19.0) 解析库确实很优秀！