#include <iostream>

#include "string.h"

extern "C" {
#include "cgo.h"
void test111();
}

std::ostream& operator<<(std::ostream& o, ParserErrorMessage_t& err);
void test_success();
void test_error();

int main() {
    test111();
    test_success();
    test_error();
}

void test_success() {
    std::cout << "test success" << std::endl;
    char input[] =
        "syntax = \"proto2\";\n"
        "message TestMessage {\n"
        "   required int32 foo = 1;\n"
        "}\n";

    ParserOption option{
        .message_type = 1,
        .require_syntax_identifier = 1,
    };
    ParserResult_t* out = ParserPBFile(input, (int)strlen(input), &option);
    if (out->errors_size) {
        for (int x = 0; x < out->errors_size; x++) {
            std::cout << out->errors[x] << std::endl;
        }
    }
    if (out->desc_size) {
        std::cout << "desc: " << out->desc << std::endl;
    }
    DeleteParserResult(out);
}

void test_error() {
    std::cout << "test error" << std::endl;
    char input[] =
        "message TestMessage {\n"
        "   required int32 foo = 1;\n"
        "}\n";

    ParserOption option{
        .message_type = 1,
        .require_syntax_identifier = 1,
    };
    ParserResult_t* out = ParserPBFile(input, (int)strlen(input), &option);
    if (out->errors_size) {
        for (int x = 0; x < out->errors_size; x++) {
            std::cout << out->errors[x] << std::endl;
        }
    }
    if (out->desc_size) {
        std::cout << "desc: " << out->desc << std::endl;
    }
    DeleteParserResult(out);
}

std::ostream& operator<<(std::ostream& o, ParserErrorMessage_t& err) {
    std::cout << "line: " << err.line << ", column: " << err.line << ", message: " << err.message << std::endl;
}

extern "C" {
void test111() {
    std::cout << "test extern C" << std::endl;
}
}