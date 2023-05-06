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
*/
import "C"

import (
	"unsafe"
)

func cgo_parse_pb(file []byte, ops ...OptionFunc) ([]byte, error) {
	if len(file) == 0 {
		return nil, ErrorMessages{{Message: `invalid proto file`}}
	}
	gooption := loadOptions(ops...)

	option := C.struct_ParserOption{}
	option.message_type = C.int(int(gooption.MessageType))
	option.require_syntax_identifier = C.int(bool2int(gooption.RequireSyntaxIdentifier))

	result := C.ParserPBFile((*C.char)(unsafe.Pointer(&file[0])), C.int(len(file)), &option)
	defer func() {
		C.DeleteParserResult(result)
	}()
	errorSize := int(result.errors_size)
	if errorSize > 0 {
		errMsg := make([]*ErrorMessage, 0, errorSize)
		for x := 0; x < errorSize; x++ {
			index := C.int(x)
			errMsg = append(errMsg, &ErrorMessage{
				Message: C.GoStringN(C.ParserResult_errors_message(result, index), C.ParserResult_errors_message_size(result, index)),
				Column:  int(C.ParserResult_errors_column(result, index)),
				Line:    int(C.ParserResult_errors_line(result, index)),
			})
		}
		return nil, ErrorMessages(errMsg)
	}
	// copy. 方便直接回收C的内存
	data := C.GoStringN(result.desc, result.desc_size)
	return UnsafeBytes(data), nil
}
