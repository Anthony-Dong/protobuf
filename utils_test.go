package protobuf

import (
	"testing"
)

func TestNewProtobufDiskSourceTree(t *testing.T) {
	tree, err := NewProtobufDiskSourceTree("internal/test/idl_example")
	if err != nil {
		t.Fatal(err)
	}
	for filename, content := range tree {
		t.Logf("file: %s, content: %d\n", filename, len(content))
	}
}

func TestUnsafeBytes(t *testing.T) {
	if UnsafeString(UnsafeBytes(`hello world`)) != `hello world` {
		t.Fatal("error")
	}
}

func TestErrors(t *testing.T) {
	messages := ErrorMessages{
		{
			Line:    1,
			Column:  1,
			Message: "line 1 message 1",
		},
		{
			Line:    2,
			Column:  2,
			Message: "line 2 message 2",
		},
	}
	t.Log(messages.Error())
}
