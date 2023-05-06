#include "pb_parser.h"

#include <google/protobuf/util/json_util.h>

namespace parser {
namespace pb {

std::ostream& operator<<(std::ostream& s, ErrorInfo& item) {
    s << "line:" << item.line << ", column: " << item.column << ", message: " << item.message;
    return s;
}

void DefaultErrorCollector::AddError(int line, google::protobuf::io::ColumnNumber column, const std::string& message) {
    this->err_info->push_back(ErrorInfo{
        .line = line,
        .column = int(column),
        .message = message,
    });
}

std::ostream& operator<<(std::ostream& s, DefaultErrorCollector& c) {
    for (int i = 0; i < int(c.err_info->size()); ++i) {
        auto item = c.err_info->at(i);
        s << item;
        if (i != 0) {
            s << std::endl;
        }
    }
    return s;
}

std::unique_ptr<std::string> MessageToBinary(google::protobuf::Message& msg) {
    using namespace google::protobuf;
    auto str = std::unique_ptr<std::string>(new std::string());
    str->reserve(64);
    io::StringOutputStream sstream(str.get());
    io::CodedOutputStream output(&sstream);
    msg.SerializePartialToCodedStream(&output);
    return str;
}

std::unique_ptr<std::string> MessageToJSON(google::protobuf::Message& msg) {
    using namespace google::protobuf;
    auto out = std::unique_ptr<std::string>(new std::string());
    out->reserve(16);
    util::MessageToJsonString(msg, out.get());
    return out;
}

std::unique_ptr<std::string> ParserPBFile(const char* file, size_t size, google::protobuf::io::ErrorCollector* error_collector, PbParserOption& option) {
    using namespace google::protobuf;
    using namespace std;

    compiler::Parser parser;
    parser.SetRequireSyntaxIdentifier(option.require_syntax_identifier);
    parser.RecordErrorsTo(error_collector);
    //    parser.SetStopAfterSyntaxIdentifier(true);
    io::ArrayInputStream raw_input(file, int(size));
    io::Tokenizer input(&raw_input, error_collector);

    FileDescriptorProto desc;
    if (!parser.Parse(&input, &desc)) {
        return std::unique_ptr<std::string>(new std::string());
    }
    switch (option.messageType) {
        case PbParserOption::Json:
            return MessageToJSON(desc);
        default:
            return MessageToBinary(desc);
    }
}

char* DupString(std::string& str) {
    char* desc_str = new char[str.length()];
    memmove(desc_str, str.c_str(), str.length());
    return desc_str;
}

#ifdef __cplusplus /* If this is a C++ compiler, use C linkage */
extern "C" {
#include "cgo.h"
#endif

struct ParserResult* ParserPBFile(const char* file, int size, struct ParserOption* option) {
    using namespace parser::pb;
    DefaultErrorCollector collector;
    PbParserOption cpp_option{
        .messageType = PbParserOption::MessageType(option->message_type),
        .require_syntax_identifier = bool(option->require_syntax_identifier),
    };
    auto desc = ParserPBFile(file, static_cast<size_t>(size), &collector, cpp_option);
    auto result = new ParserResult{};

    result->desc_size = int(desc->length());
    result->desc = DupString(*desc);
    result->errors_size = int(collector.err_info->size());
    result->errors = new ParserErrorMessage[collector.err_info->size()];

    for (int x = 0; x < int(collector.err_info->size()); x++) {
        auto item = collector.err_info->at(x);
        result->errors[x] = ParserErrorMessage{
            .line = item.line,
            .column = item.column,
            .message = parser::pb::DupString(item.message),
            .message_size = int(item.message.length()),
        };
    }
    return result;
}

void DeleteParserResult(struct ParserResult* r) {
    delete r->desc;
    for (int i = 0; i < r->errors_size; ++i) {
        delete r->errors[i].message;
    }
    delete r->errors;
    delete r;
}

#ifdef __cplusplus /* If this is a C++ compiler, end C linkage */
}
#endif

}  // namespace pb
}  // namespace parser
