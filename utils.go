package protobuf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnsafeBytes(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

func UnsafeString(bytes []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

func bool2char(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func MessageToJson(v proto.Message, pretty ...bool) string {
	ops := protojson.MarshalOptions{Multiline: len(pretty) > 0 && pretty[0]}
	return ops.Format(v)
}

func NewProtobufDiskSourceTree(dir string) (map[string][]byte, error) {
	return NewDiskSourceTree(dir, func(filename string) bool {
		return filepath.Ext(filename) == ".proto"
	})
}

func NewDiskSourceTree(dir string, filter func(filename string) bool) (map[string][]byte, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf(`filepath.Abs("%s") return err: %v`, dir, err)
	}
	files, err := lookupFiles(dir, filter)
	if err != nil {
		return nil, err
	}
	r := make(map[string][]byte, len(files))
	for _, filename := range files {
		rel, err := filepath.Rel(dir, filename)
		if err != nil {
			return nil, fmt.Errorf(`filepath.Rel("%s","%s") return err: %v`, dir, filename, err)
		}
		fileContent, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf(`ioutil.ReadFile("%s") return err: %v`, filename, err)
		}
		r[rel] = fileContent
	}
	return r, nil
}

// dir: must abs path
func lookupFiles(dir string, filter func(filename string) bool) ([]string, error) {
	files := make([]string, 0)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if filter != nil && !filter(path) {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}
