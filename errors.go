package protobuf

import (
	"fmt"
	"strings"
)

type ErrorMessage struct {
	Line    int
	Column  int
	Message string
}

type ErrorMessages []*ErrorMessage

func (e ErrorMessages) Error() string {
	result := strings.Builder{}
	result.WriteString("parse pb idl find err: ")
	for index, elem := range e {
		result.WriteString(fmt.Sprintf(`line=%d, column=%d, message=%s`, elem.Line, elem.Column, elem.Message))
		if index != len(e)-1 {
			result.WriteString("; ")
		}
	}
	return result.String()
}
