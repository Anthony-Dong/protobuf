package protobuf

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func Test_ParsePBFileDesc(t *testing.T) {
	desc, err := ParsePBFileDesc(loadIdl(t), WithRequireSyntaxIdentifier(), WithSourceCodeInfo())
	if err != nil {
		t.Fatal(err)
	}
	//writeContent(t, "internal/idl/test.proto.json", UnsafeBytes(MessageToJson(desc, true)))

	output := MessageToJson(desc, true)
	wantOutput := string(readContent(t, "internal/idl/test.proto.json"))
	assertEqualJSON(t, output, wantOutput)
}

func TestParsePBMultiFileDesc(t *testing.T) {
	desc, err := ParsePBMultiFileDesc(loadIdls(t), WithRequireSyntaxIdentifier(), WithJsonTag(), WithGoogleProtobuf(), WithSourceCodeInfo())
	if err != nil {
		t.Fatal(err)
	}
	//writeContent(t, "internal/idl/descriptor_set.json", UnsafeBytes(MessageToJson(desc, true)))

	output := MessageToJson(desc, true)
	wantOutput := string(readContent(t, "internal/idl/descriptor_set.json"))
	assertEqualJSON(t, output, wantOutput)
}

func writeContent(t testing.TB, filename string, content []byte) {
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		t.Fatal(err)
	}
}

func readContent(t testing.TB, filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func assertEqualJSON(t testing.TB, output, wantOutput string) {
	toMap := func(jj string) map[string]interface{} {
		r := make(map[string]interface{})
		if err := json.Unmarshal([]byte(jj), &r); err != nil {
			t.Fatal(err)
		}
		return r
	}
	if !reflect.DeepEqual(toMap(output), toMap(wantOutput)) {
		t.Logf("============== OUTPUT ==============\n")
		t.Log(output)
		t.Logf("============== WANT-OUTPUT ==============\n")
		t.Log(wantOutput)
		t.Fatal("============= Assert JSON ERROR =============")
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

func TestParsePBMultiFileDesc_Error(t *testing.T) {
	type args struct {
		idl *IDLConfig
		ops []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *descriptor.FileDescriptorSet
		wantErr bool
	}{
		{
			args: args{
				idl: nil,
				ops: nil,
			},
			wantErr: true,
		},
		{
			args: args{
				idl: &IDLConfig{
					Main: "xx",
					IDLs: map[string][]byte{
						"xx": {},
					},
					IncludePath: []string{},
				},
				ops: nil,
			},
			wantErr: true,
		},
		{
			args: args{
				idl: &IDLConfig{
					Main: "xx",
					IDLs: map[string][]byte{
						"xx":  []byte(`hello world`),
						"xx1": []byte(``),
					},
					IncludePath: []string{},
				},
				ops: nil,
			},
			wantErr: true,
		},
		{
			args: args{
				idl: &IDLConfig{
					Main: "xx.proto",
					IDLs: map[string][]byte{
						"xx.proto": []byte(`hello world`),
					},
					IncludePath: []string{},
				},
				ops: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePBMultiFileDesc(tt.args.idl, tt.args.ops...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePBMultiFileDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePBMultiFileDesc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
