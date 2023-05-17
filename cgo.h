#ifndef PROTOBUF_CGO_H
#define PROTOBUF_CGO_H

typedef struct ParserResult ParserResult_t;
typedef struct ParserErrorMessage ParserErrorMessage_t;
typedef struct ParserOption ParserOption_t;
typedef struct ParerFile ParerFile_t;
typedef struct ParerFiles ParerFiles_t;
typedef char BOOL;

struct ParserErrorMessage {
    int line;
    int column;
    char* message;
    int message_size;
    char* filename;
};

struct ParserOption {
    int message_type;                // 0:PB; 1:JSON
    BOOL require_syntax_identifier;  // 0:false; 1:true
    char** include_name;             // a a/b  a/b/c
    int include_name_size;
    BOOL with_source_code_info;
    BOOL with_json_tag;
    BOOL with_google_protobuf;
};

struct ParserResult {
    void* desc_pointer;            // std::string
    ParserErrorMessage_t* errors;  // if not nil will return error!
    int errors_size;
};

struct ParerFile {
    const char* filename;
    const char* content;
    int content_size;
};

struct ParerFiles {
    const char* main_filename;
    struct ParerFile* files;
    int files_size;
};

extern struct ParserResult* ParsePBFile_C(const char* file, int size, struct ParserOption* option);
extern struct ParserResult* ParseMultiPBFile_C(struct ParerFiles* files, struct ParserOption* c_option);
extern void Delete_ParserResult(struct ParserResult* r);
extern const char* ParserResult_GetDesc(struct ParserResult* r);
extern int ParserResult_GetDescSize(struct ParserResult* r);
#endif  // PROTOBUF_CGO_H