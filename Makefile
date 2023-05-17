.PHONY: all run fmt test benchmark

fmt:
	clang-format -style=file -i pb_parser.cpp pb_parser.h cgo.h
	golangci-lint run --fix -v

clean:
	rm -rf pb_parser.test mem.out

test: export CGO_ENABLED=1
test: ## go tool cover -html=coverage.out
	rm -rf internal/test/pb_gen
	go test -coverprofile=coverage.out -count=1 ./...
	cd internal/benchmark && go test -count=1 ./...
	go run internal/example/main.go > /dev/null 2>&1 || exit 1

benchmark_cgo: export CGO_ENABLED=1
benchmark_cgo: ## -gcflags="-N -l -m"
	go test -v -run=none -bench=Benchmark -memprofile mem.out  -benchmem  -count=5 .

benchmark: export CGO_ENABLED=1
benchmark: benchmark_cgo
	@cd internal/benchmark && \
	go test -v -run=none -bench=Benchmark -memprofile mem_bench.out  -benchmem  -count=5 .

gogoprotobuf:
