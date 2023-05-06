#include "cgo.h"
#include "stdio.h"
#include "string.h"

int main() {
    char input[] =
        "syntax = \"proto2\";"
        "message TestMessage {\n"
        "  required int32 foo = 1;\n"
        "}\n";
    ParserOption_t t;
    t.message_type = 1;
    t.require_syntax_identifier = 1;
    ParserResult_t* out = ParserPBFile(input, (int)strlen(input), &t);
    printf("%s\n", out->desc);
    printf("%d\n", out->desc_size);
    printf("%p\n", out->errors);
    printf("%d\n", out->errors_size);

    DeleteParserResult(out);
}