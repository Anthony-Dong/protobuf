#ifndef PROTOBUF_PB_PARSER_H
#define PROTOBUF_PB_PARSER_H

#include <google/protobuf/compiler/parser.h>

namespace parser {
namespace pb {

class ErrorInfo;
class PbParserOption;
class DefaultErrorCollector;

std::unique_ptr<std::string> MessageToBinary(google::protobuf::Message& msg);
std::unique_ptr<std::string> MessageToJSON(google::protobuf::Message& msg);
std::unique_ptr<std::string> ParserPBFile(const char* file, size_t size, google::protobuf::io::ErrorCollector* error_collector, PbParserOption& option);
char* DupString(std::string& str);

class ErrorInfo {
   public:
    int line;
    int column;
    std::string message;

   public:
    friend std::ostream& operator<<(std::ostream& s, ErrorInfo& item);
};

class DefaultErrorCollector : public google::protobuf::io::ErrorCollector {
   public:
    DefaultErrorCollector() noexcept {
        this->err_info.reset(new std::vector<ErrorInfo>());
    };
    ~DefaultErrorCollector() override = default;
    DefaultErrorCollector(const DefaultErrorCollector&) = delete;
    void operator=(const DefaultErrorCollector&) = delete;

    void AddError(int line, google::protobuf::io::ColumnNumber column, const std::string& message) override;
    void AddWarning(int line, google::protobuf::io::ColumnNumber column, const std::string& message) override{};
    friend std::ostream& operator<<(std::ostream& s, DefaultErrorCollector& c);

   public:
    std::unique_ptr<std::vector<ErrorInfo>> err_info;
};

class PbParserOption {
   public:
    PbParserOption() = default;
    PbParserOption(const PbParserOption&) = delete;
    void operator=(const PbParserOption&) = delete;

    enum MessageType : int {
        PB,
        Json,
    };

   public:
    MessageType messageType;
    bool require_syntax_identifier;
};
}  // namespace pb
}  // namespace parser

#endif  // PROTOBUF_PB_PARSER_H
