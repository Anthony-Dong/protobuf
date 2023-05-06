#ifndef PROTOBUF_CGO_H
#define PROTOBUF_CGO_H

typedef struct ParserResult ParserResult_t;
typedef struct ParserErrorMessage ParserErrorMessage_t;
typedef struct ParserOption ParserOption_t;

struct ParserErrorMessage {
    int line;
    int column;
    char* message;
    int message_size;
};

struct ParserOption {
    int message_type;               // 0:PB; 1:JSON
    int require_syntax_identifier;  // 0:false; 1:true
};

struct ParserResult {
    char* desc;
    int desc_size;
    ParserErrorMessage_t* errors;  // if not nil will return error!
    int errors_size;
};

extern struct ParserResult* ParserPBFile(const char* file, int size, struct ParserOption* option);
extern void DeleteParserResult(struct ParserResult* r);
#endif  // PROTOBUF_CGO_H