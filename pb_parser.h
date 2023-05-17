#ifndef PROTOBUF_PB_PARSER_H
#define PROTOBUF_PB_PARSER_H

#include <google/protobuf/compiler/importer.h>
#include <google/protobuf/compiler/parser.h>
#include <google/protobuf/io/zero_copy_stream_impl_lite.h>

#include <unordered_map>

namespace parser {
namespace pb {

class PbParserOption;
class WrappedErrorCollector;
class MemorySourceTree;
struct StringView;
typedef WrappedErrorCollector Errors;

typedef std::unordered_map<std::string, StringView> StringMap;

std::unique_ptr<std::string> ParsePBFile(const char* file, size_t size, google::protobuf::io::ErrorCollector* error_collector, PbParserOption* option);
std::unique_ptr<std::string> ParseMultiPBFile(const std::string& main, const StringMap& treeMap, Errors* errors, PbParserOption* option);

struct StringView {
    const char* data_;
    size_t size_;
};

class WrappedErrorCollector : public google::protobuf::io::ErrorCollector, public google::protobuf::compiler::MultiFileErrorCollector, public google::protobuf::DescriptorPool::ErrorCollector {
   public:
    struct ErrorInfo {
        std::string filename;
        std::string elemName;
        int line;
        int column;
        std::string message;
    };
    typedef std::vector<ErrorInfo> ErrorInfoList;

   public:
    WrappedErrorCollector() = default;
    ~WrappedErrorCollector() override = default;
    WrappedErrorCollector(const WrappedErrorCollector&) = delete;
    void operator=(const WrappedErrorCollector&) = delete;

    // Indicates that there was an error in the input at the given line and
    // column numbers.  The numbers are zero-based, so you may want to add
    // 1 to each before printing them.
    void AddError(int line, google::protobuf::io::ColumnNumber column, const std::string& message) override;

    // Line and column numbers are zero-based.  A line number of -1 indicates
    // an error with the entire file (e.g. "not found").
    void AddError(const std::string& filename, int line, int column, const std::string& message) override;

    // Reports an error in the FileDescriptorProto. Use this function if the
    // problem occurred should interrupt building the FileDescriptorProto.
    virtual void AddError(
        const std::string& filename,                  // File name in which the error occurred.
        const std::string& element_name,              // Full name of the erroneous element.
        const google::protobuf::Message* descriptor,  // Descriptor of the erroneous element.
        ErrorLocation location,                       // One of the location constants, above.
        const std::string& message                    // Human-readable error message.
        ) override;

   public:
    ErrorInfoList errors;
};

class PbParserOption {
   public:
    enum MessageType : int {
        PB,
        Json,
    };

   public:
    PbParserOption() : message_type_(MessageType::PB), require_syntax_identifier_(false), with_source_code_info_(false){};
    PbParserOption(const PbParserOption&) = delete;
    void operator=(const PbParserOption&) = delete;

   public:
    MessageType message_type_;
    bool require_syntax_identifier_;
    std::vector<std::string> include_path_;
    bool with_source_code_info_;
    bool with_json_tag_;
    bool with_google_protobuf_;
};

// A dummy implementation of SourceTree backed by a simple map.
class MemorySourceTree : public google::protobuf::compiler::SourceTree, public google::protobuf::DescriptorDatabase {
   public:
    MemorySourceTree() = default;
    MemorySourceTree(const StringMap& files) : files_(files){};
    ~MemorySourceTree() override = default;
    void AddFile(const std::string& filename, const char* contents, size_t content_size) {
        this->files_[filename] = {
            .data_ = contents,
            .size_ = content_size,
        };
    }
    void SetErrorCollector(google::protobuf::compiler::MultiFileErrorCollector* error_collector) {
        this->error_collector_ = error_collector;
    }
    void SetPbParserOption(PbParserOption* option) {
        this->parser_option_ = option;
    }
    // implements SourceTree -----------------------------------
    google::protobuf::io::ZeroCopyInputStream* Open(const std::string& filename) override;
    std::string GetLastErrorMessage() override {
        return "File not found.";
    };

    // implements DescriptorDatabase -----------------------------------
    bool FindFileByName(const std::string& filename,
                        google::protobuf::FileDescriptorProto* output) override;
    bool FindFileContainingSymbol(const std::string& symbol_name,
                                  google::protobuf::FileDescriptorProto* output) override;
    bool FindFileContainingExtension(const std::string& containing_type,
                                     int field_number,
                                     google::protobuf::FileDescriptorProto* output) override;

   private:
    StringMap files_;
    google::protobuf::compiler::MultiFileErrorCollector* error_collector_;
    PbParserOption* parser_option_;
    GOOGLE_DISALLOW_EVIL_CONSTRUCTORS(MemorySourceTree);

   private:
    class FileErrorCollector : public google::protobuf::io::ErrorCollector {
       public:
        FileErrorCollector() = delete;
        FileErrorCollector(std::string filename, google::protobuf::compiler::MultiFileErrorCollector* error_collector) : filename_(filename), error_collector_(error_collector), had_errors_(false){};
        void AddError(int line, google::protobuf::io::ColumnNumber column,
                      const std::string& message) override {
            if (this->error_collector_ != nullptr) {
                this->error_collector_->AddError(this->filename_, line, column, message);
            }
            this->had_errors_ = true;
        };
        bool HadErrors() {
            return this->had_errors_;
        }

       private:
        std::string filename_;
        google::protobuf::compiler::MultiFileErrorCollector* error_collector_;
        bool had_errors_;
        GOOGLE_DISALLOW_EVIL_CONSTRUCTORS(FileErrorCollector);
    };
};

}  // namespace pb
}  // namespace parser

#endif  // PROTOBUF_PB_PARSER_H
