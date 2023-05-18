# Protobuf Parser

## 能力

1. 基于 libprotobuf 3.19.0 解析 proto 文件，解决一些开源的Go的Protobuf解析库和官方的库不一致的问题！

2. 支持的环境
- Darwin amd64
- Linux amd64

> 有其他环境需求可自行提交 MR，只要本地跑通单元测试即可！


## 使用

1. 下载依赖

```shell
CGO_ENABLED=1 go get -v github.com/anthony-dong/protobuf
```

2. 如何使用: [main.go](internal/example/main.go)

```go
package main

import (
	"fmt"

	"github.com/anthony-dong/protobuf"
)

func main() {
	tree, err := protobuf.NewProtobufDiskSourceTree("internal/test/idl_example")
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

// 运行: CGO_ENABLED=1 go run main.go
```

## 性能

压测文件

1. 解析PB文件性能

   >  压测文件: [cgo_test.go](./cgo_test.go)

```shell
go test -v -run=none -bench=Benchmark -memprofile mem.out  -benchmem  -count=5 .
goos: linux
goarch: amd64
pkg: github.com/anthony-dong/protobuf
cpu: Intel(R) Xeon(R) Platinum 8260 CPU @ 2.40GHz
Benchmark_cgo_parse_pb_pb
Benchmark_cgo_parse_pb_pb-8     	   29761	     40551 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   28610	     41515 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   29256	     41575 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   29144	     41192 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_pb-8     	   30019	     40678 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_json
Benchmark_cgo_parse_pb_json-8   	    9913	    114188 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_json-8   	    9848	    114095 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_json-8   	    9326	    115046 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_json-8   	   10000	    116789 ns/op	      16 B/op	       1 allocs/op
Benchmark_cgo_parse_pb_json-8   	    9115	    117856 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/anthony-dong/protobuf	13.772s
```

2. 对比 [github.com/jhump/protoreflect@v1.8.2](https://github.com/jhump/protoreflect/tree/v1.8.2 ) 解析库

> 原因是我司用的是 v1.8.2 版本, 高版本兼容性检测会高一些!
>
> 压测文件:
>
> -  [benchmark_test.go](internal/benchmark/benchmark_test.go) 
> -  [benchmark_multi_test.go](internal/benchmark/benchmark_multi_test.go)

```shell
goos: linux
goarch: amd64
pkg: github.com/anthony-dong/protobuf/internal/benchmark
cpu: Intel(R) Xeon(R) Platinum 8260 CPU @ 2.40GHz
Benchmark_ParsePBMultiFileDesc_Cgo
Benchmark_ParsePBMultiFileDesc_Cgo-8     	     942	   1207139 ns/op	   60937 B/op	    1976 allocs/op
Benchmark_ParsePBMultiFileDesc_Cgo-8     	     996	   1216023 ns/op	   60939 B/op	    1976 allocs/op
Benchmark_ParsePBMultiFileDesc_Cgo-8     	     958	   1223891 ns/op	   60936 B/op	    1976 allocs/op
Benchmark_ParsePBMultiFileDesc_Cgo-8     	     962	   1245004 ns/op	   60936 B/op	    1976 allocs/op
Benchmark_ParsePBMultiFileDesc_Cgo-8     	     986	   1225556 ns/op	   60936 B/op	    1976 allocs/op
Benchmark_ParsePBMultiFileDesc_Jhump
Benchmark_ParsePBMultiFileDesc_Jhump-8   	     508	   2404250 ns/op	  762870 B/op	   11038 allocs/op
Benchmark_ParsePBMultiFileDesc_Jhump-8   	     468	   2414933 ns/op	  762874 B/op	   11038 allocs/op
Benchmark_ParsePBMultiFileDesc_Jhump-8   	     494	   2404417 ns/op	  762828 B/op	   11038 allocs/op
Benchmark_ParsePBMultiFileDesc_Jhump-8   	     498	   2343901 ns/op	  762930 B/op	   11038 allocs/op
Benchmark_ParsePBMultiFileDesc_Jhump-8   	     494	   2412528 ns/op	  762977 B/op	   11038 allocs/op
Benchmark_ParsePBFileDesc_Cgo
Benchmark_ParsePBFileDesc_Cgo-8          	   22116	     54090 ns/op	    5024 B/op	     145 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8          	   22134	     54031 ns/op	    5024 B/op	     145 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8          	   21500	     56155 ns/op	    5024 B/op	     145 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8          	   20539	     57158 ns/op	    5024 B/op	     145 allocs/op
Benchmark_ParsePBFileDesc_Cgo-8          	   21400	     55135 ns/op	    5024 B/op	     145 allocs/op
Benchmark_ParsePBFileDesc_Jhump
Benchmark_ParsePBFileDesc_Jhump-8        	    3837	    307914 ns/op	  107032 B/op	    1735 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8        	    3874	    311001 ns/op	  107036 B/op	    1735 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8        	    3578	    319180 ns/op	  107037 B/op	    1735 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8        	    3895	    308775 ns/op	  107035 B/op	    1735 allocs/op
Benchmark_ParsePBFileDesc_Jhump-8        	    3873	    313726 ns/op	  107033 B/op	    1735 allocs/op
PASS
ok  	github.com/anthony-dong/protobuf/internal/benchmark	28.682s
```

3. 结论：

- 解析单文件实际上由于C++极致的性能，也就是说冗余了一次序列化+反序列化，性能也能是5-6倍！
- 解析多文件由于大量的时间浪费在序列化和反序列化上，导致性能裂化较为严重，性能只能到2倍左右！
- 其次就是内存压力的大大降低！


## Todo

- Go 官方的 Protobuf 序列化库实际上是用的反射去实现的，可以通过代码生成工具实现硬编码解析，性能会再提高很多！
- 寻找更多的突破口，降低序列化和反序列化的开销，例如直接做C++ -> Go的数据绑定, 这个可能代码工作量会大一些，但是可以避免无效的序列化，这个我争取在下个版本release了！
- CGO中存在大量的数据转换，是不是有一些代码生成辅助工具帮忙实现呢？但是GO和C其实还好，但是GO与C++就比较恶心了，需要保存一个`void*`指针去做中转，因为C++基本上不可能只有POD类型，大多数都是复杂类型！
