#include <iostream>

#include "string.h"

extern "C" {
#include "cgo.h"
}

#include "pb_parser.h"

const char* api_proto() {
    return "syntax = \"proto2\";\n"
           "package api;\n"
           "option go_package = \"github.com/anthony-dong/go-tool/internal/example/protobuf/idl_example/pb_gen/api\";\n"
           "\n"
           "import \"google/protobuf/descriptor.proto\";\n"
           "\n"
           "extend google.protobuf.FileOptions{\n"
           "  optional string android_package = 1001;\n"
           "}\n"
           "\n"
           "extend google.protobuf.FieldOptions {\n"
           "  optional HttpSourceType source = 50101; // 来自于http 请求的哪个部位\n"
           "  optional string key = 50102; // http 请求的header 还是哪\n"
           "  optional bool unbox = 50103; // 是否平铺开结构体，除body外，默认处理第一层\n"
           "}\n"
           "\n"
           "enum HttpSourceType {\n"
           "  Query = 1;\n"
           "  Body = 2;\n"
           "  Header = 3;\n"
           "}\n"
           "\n"
           "extend google.protobuf.MethodOptions {\n"
           "  optional HttpMethodType method = 50201; // http method\n"
           "  optional string path = 50202; // http path\n"
           "}\n"
           "\n"
           "enum HttpMethodType{\n"
           "  GET = 1;\n"
           "  POST = 2;\n"
           "  PUT = 3;\n"
           "}";
}

const char* common_proto() {
    return "syntax = \"proto2\";\n"
           "package im.commons;\n"
           "option go_package = \"github.com/anthony-dong/go-tool/internal/example/protobuf/idl_example/pb_gen/commons\";\n"
           "\n"
           "import \"api.proto\";\n"
           "\n"
           "message UserInfo {\n"
           "  optional int64 Id = 1 [(api.source) = Header, (api.key) = 'X-Biz-UserId'];\n"
           "  optional string Name = 2 [(api.source) = Header, (api.key) = 'X-Biz-UserName'];\n"
           "}";
}

const char* service_proto() {
    return "syntax = \"proto2\";\n"
           "\n"
           "package im.service;\n"
           "\n"
           "import \"api.proto\";\n"
           "import \"im/common.proto\";\n"
           "\n"
           "option go_package = \"github.com/anthony-dong/go-tool/internal/example/protobuf/idl_example/pb_gen/service\";\n"
           "\n"
           "// 自定义 安卓package!\n"
           "option (api.android_package) = \"github.com.anthony_dong.go_tool.service\";\n"
           "\n"
           "service ImService {\n"
           "  rpc GetMessage (im.commons.UserInfo) returns (im.commons.UserInfo){\n"
           "    option (api.method) = GET;\n"
           "    option (api.path) = '/api/v1/im/query';\n"
           "  };\n"
           "}";
}

void test_cpp() {
    using namespace parser::pb;
    StringMap tree;

    tree["desc/api.proto"] = {
        .data_ = api_proto(),
        .size_ = strlen(api_proto()),
    };
    tree["im/common.proto"] = {
        .data_ = common_proto(),
        .size_ = strlen(common_proto()),
    };
    tree["service.proto"] = {
        .data_ = service_proto(),
        .size_ = strlen(service_proto()),
    };

    WrappedErrorCollector collector;
    PbParserOption option;
    option.message_type_ = PbParserOption::MessageType::Json;
    option.include_path_.push_back("");
    option.include_path_.push_back("desc");
    option.require_syntax_identifier_ = true;
    option.with_source_code_info_ = true;
    auto result = ParserMultiPBFile("service.proto", tree, &collector, &option);
    if (result != nullptr) {
        std::cout << *result << std::endl;
    }
    for (const auto& err : collector.errors) {
        std::cout << err.filename << ": " << err.message << std::endl;
    }
}

void test_c() {
    ParerFile_t t[3];
    t[0].filename = "desc/api.proto";
    t[0].content = api_proto();
    t[0].content_size = strlen(api_proto());

    t[1].filename = "im/common.proto";
    t[1].content = common_proto();
    t[1].content_size = strlen(common_proto());

    t[2].filename = "service.proto";
    t[2].content = service_proto();
    t[2].content_size = strlen(service_proto());

    ParerFiles files{
        .main_filename = "service.proto",
        .files = t,
        .files_size = 3,
    };
    char* include_path[2];
    include_path[0] = "";
    include_path[1] = "desc";
    ParserOption option{
        .message_type = 1,
        .include_name = include_path,
        .include_name_size = 2,
        .with_source_code_info = 1,
    };
    auto result = ParseMultiPBFile_C(&files, &option);
    auto desc = ParserResult_GetDesc(result);
    if (desc != nullptr) {
        std::cout << desc << std::endl;
    }
    Delete_ParserResult(result);
}

int main() {
    test_cpp();
    test_c();
}