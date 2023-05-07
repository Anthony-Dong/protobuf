package protobuf

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func Test_ParsePBFileDesc(t *testing.T) {
	desc, err := ParsePBFileDesc(loadIdl(t), WithRequireSyntaxIdentifier())
	if err != nil {
		t.Fatal(err)
	}
	jsonMsg := MessageToJson(desc, true)

	//if err := ioutil.WriteFile("internal/test/api.proto.json", []byte(jsonMsg), 0644); err != nil {
	//	t.Fatal(err)
	//}

	assertFiles, err := ioutil.ReadFile("internal/test/api.proto.json")
	if err != nil {
		t.Fatal(err)
	}

	toMap := func(jj string) map[string]interface{} {
		r := make(map[string]interface{})
		if err := json.Unmarshal([]byte(jj), &r); err != nil {
			t.Fatal(err)
		}
		return r
	}
	if !reflect.DeepEqual(toMap(string(assertFiles)), toMap(jsonMsg)) {
		t.Fatal("assert pb desc find err")
	}
}

func Test_ParsePBFileDesc_Error(t *testing.T) {
	type args struct {
		file []byte
		ops  []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *descriptor.FileDescriptorProto
		wantErr bool
	}{
		{
			name: "null",
			args: args{
				file: nil,
			},
			wantErr: true,
		},
		{
			name: "null",
			args: args{
				file: []byte(`
hello world
`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePBFileDesc(tt.args.file, tt.args.ops...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePBFileDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePBFileDesc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
