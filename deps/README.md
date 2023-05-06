# libprotobuf 3.19.0.0

>  注意Go CGO最好是把依赖的静态库直接放到你的仓库，这么用户就不用自己本地构建了！

## 下载依赖

```shell
# 1. 下载protoc
wget https://github.com/protocolbuffers/protobuf/archive/refs/tags/v3.19.0.tar.gz -O- | tar -zxvf -

# 进入目录
cd protobuf-3.19.0

# 2. 进入release目录
cd cmake && mkdir -p release && cd release


# 3. 构建cmake, 这里是静态依赖
cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release  -DCMAKE_CXX_STANDARD=11  -Dprotobuf_BUILD_TESTS=OFF -DCMAKE_INSTALL_PREFIX=/usr/local  ../.

# 4. 构建
make

# 5. 安装，注意: 这里不需要install 到 /usr/local 下, 如果你需要的话执行，sudo make install
```

### 注意：

`CGO`使用动态链接库的话，很容易找不到动态库在哪里，所以会报以下错误: 

```shell
go test -coverprofile=coverage.out -count=1 ./...
go: downloading google.golang.org/protobuf v1.30.0
/tmp/go-build006205204/b001/protobuf.test: error while loading shared libraries: libprotobuf.so.3.19.0.0: cannot open shared object file: No such file or directory
FAIL	github.com/anthony-dong/protobuf	0.001s
?   	github.com/anthony-dong/protobuf/deps/darwin_x86_64	[no test files]
?   	github.com/anthony-dong/protobuf/deps/linux_x86_64	[no test files]
?   	github.com/anthony-dong/protobuf/internal/example	[no test files]
FAIL
make: *** [Makefile:12: test] Error 1
Error: Process completed with exit code 2.
```

如何解决动态库的问题了，这时候需要你自己导出路径，`${YOUR_LIB_DIR}` 改成你的 `-L dir`

```shell
export LD_LIBRARY_PATH=${YOUR_LIB_DIR}:$LD_LIBRARY_PATH
```

**所以一般都是用的静态库，可以避免这些坑！！！**

## 复制到此项目

文件目录格式为 `os_{amd/arm}`

- amd:  `x86_64`
- arm: `arm64`

```shell
cd ${DIR}
cp /home/fanhaodong.516/data/protobuf-3.19.0/cmake/release/libprotobuf.a  .
```

注意: 只需要拷贝 `libprotobuf` 文件即可，仅支持64位系统以及C++11 ！
