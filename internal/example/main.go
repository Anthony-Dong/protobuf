package main

import (
	"log"

	"github.com/anthony-dong/protobuf"
)

func main() {
	file := []byte(`
syntax = "proto2";
package idl.model;
message Person {
  optional string name = 1;
  optional int32 id = 2;
  optional string email = 3;
  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
  }
  message PhoneNumber {
    optional string number = 1;
    optional PhoneType type = 2 [default = HOME];
  }
  repeated PhoneNumber phones = 4;
  map<string, Person> map_person = 5;
  optional bool status = 6;
}
`)
	desc, err := protobuf.ParsePBFileDesc(file, protobuf.WithRequireSyntaxIdentifier())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(protobuf.MessageToJson(desc))
}
