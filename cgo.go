package protobuf

/*
#cgo CXXFLAGS: -std=c++11 -Wall -I${SRCDIR}/deps/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/deps/darwin_x86_64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/deps/linux_x86_64
#cgo LDFLAGS: -lprotobuf
#include <stdlib.h>
#include <cgo.h>

int ParserResult_errors_line(ParserResult_t* r, int index) {
    return r->errors[index].line;
}

int ParserResult_errors_column(ParserResult_t* r, int index) {
    return r->errors[index].column;
}

const char* ParserResult_errors_message(ParserResult_t* r, int index) {
    return r->errors[index].message;
}

int ParserResult_errors_message_size(ParserResult_t* r, int index){
    return r->errors[index].message_size;
}

char* ParserResult_errors_filename(ParserResult_t* r, int index){
    return r->errors[index].filename;
}

ParserOption_t* New_ParserOption() {
    return (ParserOption_t*)calloc(1, sizeof(ParserOption_t));
}

void ParserOption_init_include_path(ParserOption_t* r, int size){
    r->include_name = (char**)malloc(sizeof(char*)*size);
	r->include_name_size = size;
}

void ParserOption_add_include_path(ParserOption_t* r, int index, char* include){
    r->include_name[index] = include;
}

void Delete_ParserOption(ParserOption_t* r) {
	free(r->include_name);
	free(r);
}

ParerFiles_t* New_ParerFiles(char* main_filename, int file_size){
	ParerFiles_t* r = (ParerFiles_t*)malloc(sizeof(ParerFiles_t));
	r->files = (ParerFile_t*)malloc(sizeof(ParerFile_t) * file_size);
	r->files_size = file_size;
	r->main_filename = main_filename;
	return r;
}

void Delete_ParerFiles(ParerFiles_t* r) {
	free(r->files);
	free(r);
}

void ParerFiles_Set_File(ParerFiles_t* r,int index, char* filename, char* content, int content_size){
	r->files[index].filename = filename;
	r->files[index].content = content;
	r->files[index].content_size = content_size;
}

*/
import "C"

import (
	"reflect"
	"strings"
	"unsafe"
)

const cStrEnd = string('\u0000')

func unsafeCString(str string) *C.char {
	// C语言的字符串是以\u0000 结尾的，所以这里注意了. 需要手动加一个结尾符号
	if index := strings.IndexByte(str, '\u0000'); index == -1 {
		str = str + cStrEnd
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*C.char)(unsafe.Pointer(header.Data))
}

func safeCString(str string) *C.char {
	return C.CString(str)
}

func unsafeCBytes(str []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&str[0]))
}

func unsafeGoBytes(arr *C.char, arrSize C.int) []byte {
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(arr)),
		Len:  int(arrSize),
		Cap:  int(arrSize),
	}
	return *(*[]byte)(unsafe.Pointer(&header))
}

func newParserOption(ops *Option, includes []string) (*C.ParserOption_t, func()) {
	if ops == nil {
		ops = &Option{}
	}
	option := C.New_ParserOption() // GO用指针会逃逸到堆上，所以直接malloc堆即可
	option.message_type = C.int(int(ops.MessageType))
	option.require_syntax_identifier = C.char(bool2char(ops.RequireSyntaxIdentifier))
	option.with_source_code_info = C.char(bool2char(ops.WithSourceCodeInfo))
	option.with_json_tag = C.char(bool2char(ops.WithJsonTag))
	option.with_google_protobuf = C.char(bool2char(ops.WithGoogleProtobuf))
	free := func() {
		C.Delete_ParserOption(option)
	}
	if len(includes) == 0 {
		return option, free
	}
	C.ParserOption_init_include_path(option, C.int(len(includes)))
	for index, include := range includes {
		C.ParserOption_add_include_path(option, C.int(index), unsafeCString(include))
	}
	return option, free
}

func newParerFiles(idlConfig *IDLConfig) (*C.ParerFiles_t, func()) {
	files := C.New_ParerFiles(unsafeCString(idlConfig.Main), C.int(len(idlConfig.IDLs)))
	index := 0
	for filename, content := range idlConfig.IDLs {
		C.ParerFiles_Set_File(files, C.int(index), unsafeCString(filename), unsafeCBytes(content), C.int(len(content)))
		index = index + 1
	}
	return files, func() {
		C.Delete_ParerFiles(files)
	}
}

func getErrorMessage(result *C.ParserResult_t) ErrorMessages {
	errorSize := int(result.errors_size)
	if errorSize > 0 {
		errMsg := make([]*ErrorMessage, 0, errorSize)
		for x := 0; x < errorSize; x++ {
			index := C.int(x)
			errMsg = append(errMsg, &ErrorMessage{
				Message:  C.GoStringN(C.ParserResult_errors_message(result, index), C.ParserResult_errors_message_size(result, index)),
				Column:   int(C.ParserResult_errors_column(result, index)),
				Line:     int(C.ParserResult_errors_line(result, index)),
				Filename: C.GoString(C.ParserResult_errors_filename(result, index)),
			})
		}
		return errMsg
	}
	return nil
}

func cgo_parse_pb(file []byte, ops *Option, handler func(desc []byte) error) error {
	if len(file) == 0 {
		return ErrorMessages{{Message: `invalid proto file`}}
	}
	option, free := newParserOption(ops, nil)
	defer free()
	result := C.ParsePBFile_C(unsafeCBytes(file), C.int(len(file)), option)
	defer func() {
		C.Delete_ParserResult(result)
	}()
	if err := getErrorMessage(result); err != nil {
		return err
	}
	return handler(unsafeGoBytes(C.ParserResult_GetDesc(result), C.ParserResult_GetDescSize(result)))
}

func cgo_parse_multi_pb(idlConfig *IDLConfig, ops *Option, handler func(desc []byte) error) error {
	if idlConfig == nil || idlConfig.Main == "" || len(idlConfig.IDLs) == 0 || len(idlConfig.IDLs[idlConfig.Main]) == 0 {
		return ErrorMessages{{Message: `invalid proto file`}}
	}
	for _, idl := range idlConfig.IDLs {
		if len(idl) == 0 {
			return ErrorMessages{{Message: `invalid proto file`}}
		}
	}
	option, free := newParserOption(ops, idlConfig.IncludePath)
	defer free()
	files, freeFiles := newParerFiles(idlConfig)
	defer freeFiles()
	result := C.ParseMultiPBFile_C(files, option)
	defer func() {
		C.Delete_ParserResult(result)
	}()
	if err := getErrorMessage(result); err != nil {
		return err
	}
	if handler == nil {
		return nil
	}
	return handler(unsafeGoBytes(C.ParserResult_GetDesc(result), C.ParserResult_GetDescSize(result)))
}
