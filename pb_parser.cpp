
#include "pb_parser.h"

#include <google/protobuf/compiler/plugin.pb.h>
#include <google/protobuf/util/json_util.h>

#include "pb_include.h"

namespace parser {
namespace pb {

using namespace google::protobuf;

void WrappedErrorCollector::AddError(int line, io::ColumnNumber column, const std::string& message) {
    this->errors.push_back(ErrorInfo{
        .filename = "",
        .elemName = "",
        .line = line,
        .column = int(column),
        .message = message,
    });
}

void WrappedErrorCollector::AddError(const std::string& filename, int line, int column, const std::string& message) {
    this->errors.push_back(ErrorInfo{
        .filename = filename,
        .elemName = {},
        .line = line,
        .column = int(column),
        .message = message,
    });
}

void WrappedErrorCollector::AddError(const string& filename, const string& element_name, const Message* descriptor, ErrorLocation location, const string& message) {
    this->errors.push_back(ErrorInfo{
        .filename = filename,
        .elemName = element_name,
        .line = -1,
        .column = 0,
        .message = message,
    });
}

io::ZeroCopyInputStream* MemorySourceTree::Open(const std::string& filename) {
    if (this->parser_option_ == nullptr || this->parser_option_->include_path_.empty()) {
        auto file = this->files_.find(filename);
        if (file == this->files_.end()) {
            return nullptr;
        }
        return new io::ArrayInputStream(file->second.data_, int(file->second.size_));
    }
    for (const auto& include_path : this->parser_option_->include_path_) {
        decltype(this->files_.end()) file;
        if (include_path == "" || include_path == ".") {
            file = this->files_.find(filename);
        } else {
            file = this->files_.find(include_path + "/" + filename);
        }
        if (file == this->files_.end()) {
            continue;
        }
        return new io::ArrayInputStream(file->second.data_, int(file->second.size_));
    }
    return nullptr;
}

bool findGoogleProtobuf(const ::std::string& filename, google::protobuf::FileDescriptorProto* output);

bool MemorySourceTree::FindFileByName(const std::string& filename, google::protobuf::FileDescriptorProto* output) {
    std::unique_ptr<io::ZeroCopyInputStream> input_stream(this->Open(filename));
    if (input_stream == nullptr) {
        if (findGoogleProtobuf(filename, output)) {
            return true;
        }
        if (this->error_collector_ != nullptr) {
            this->error_collector_->AddError(filename, -1, 0, this->GetLastErrorMessage());
        }
        return false;
    }
    compiler::Parser parser;
    MemorySourceTree::FileErrorCollector file_error_collector(filename, this->error_collector_);
    parser.RecordErrorsTo(&file_error_collector);
    if (this->parser_option_ != nullptr) {
        parser.SetRequireSyntaxIdentifier(this->parser_option_->require_syntax_identifier_);
    }
    io::Tokenizer tokenizer(input_stream.get(), &file_error_collector);
    output->set_name(filename);
    return (parser.Parse(&tokenizer, output) && !file_error_collector.HadErrors());
}

bool MemorySourceTree::FindFileContainingSymbol(const std::string& symbol_name, google::protobuf::FileDescriptorProto* output) {
    return false;
}

bool MemorySourceTree::FindFileContainingExtension(const std::string& containing_type, int field_number, google::protobuf::FileDescriptorProto* output) {
    return false;
}

inline std::unique_ptr<std::string> messageToBinary(Message& msg);
inline std::unique_ptr<std::string> messageToJSON(Message& msg);

std::unique_ptr<std::string> ParsePBFile(const char* file, size_t size, io::ErrorCollector* error_collector, PbParserOption* option) {
    using namespace std;
    compiler::Parser parser;
    parser.SetRequireSyntaxIdentifier(option->require_syntax_identifier_);
    parser.RecordErrorsTo(error_collector);
    io::ArrayInputStream raw_input(file, int(size));
    io::Tokenizer input(&raw_input, error_collector);
    FileDescriptorProto desc;
    if (!parser.Parse(&input, &desc)) {
        return std::unique_ptr<std::string>();
    }
    if (!option->with_source_code_info_) {
        desc.clear_source_code_info();
    }
    switch (option->message_type_) {
        case PbParserOption::Json:
            return messageToJSON(desc);
        default:
            return messageToBinary(desc);
    }
}

typedef WrappedErrorCollector Errors;
bool loadFileDescriptorSet(FileDescriptorSet* files, const FileDescriptor* main, Errors* errors, std::unordered_set<std::string>* walk_set, PbParserOption* option);

std::unique_ptr<std::string> ParseMultiPBFile(const std::string& main, const StringMap& treeMap, Errors* errors, PbParserOption* option) {
    MemorySourceTree sourceTree(treeMap);
    sourceTree.SetErrorCollector(errors);
    sourceTree.SetPbParserOption(option);

    DescriptorPool descriptorPool(&sourceTree, errors);
    auto mainDesc = descriptorPool.FindFileByName(main);
    if (mainDesc == nullptr) {
        if (errors->errors.empty()) {
            errors->AddError(main, -1, 0, "Not Found File");
        }
        return std::unique_ptr<std::string>();
    }
    std::unordered_set<std::string> walk_set;
    FileDescriptorSet fileSet;
    if (!loadFileDescriptorSet(&fileSet, mainDesc, errors, &walk_set, option)) {
        return std::unique_ptr<std::string>();
    }
    PbParserOption::MessageType message_type = PbParserOption::MessageType::PB;
    if (option != nullptr) {
        message_type = option->message_type_;
    }
    switch (message_type) {
        case PbParserOption::MessageType::Json:
            return messageToJSON(fileSet);
        default:
            return messageToBinary(fileSet);
    }
}

inline char* dupString(std::string& str) {
    char* desc_str = new char[str.length()];
    if (str.length() == 0) {
        return desc_str;
    }
    memmove(desc_str, str.c_str(), str.length());
    return desc_str;
}

inline std::unique_ptr<std::string> messageToBinary(Message& msg) {
    auto str = std::unique_ptr<std::string>(new std::string());
    str->reserve(64);
    io::StringOutputStream sstream(str.get());
    io::CodedOutputStream output(&sstream);
    msg.SerializePartialToCodedStream(&output);
    return str;
}

inline std::unique_ptr<std::string> messageToJSON(Message& msg) {
    auto out = std::unique_ptr<std::string>(new std::string());
    out->reserve(64 * 2);
    util::MessageToJsonString(msg, out.get());
    return out;
}

bool findGoogleProtobuf(const ::std::string& filename, google::protobuf::FileDescriptorProto* output) {
    typedef ::std::unordered_map<::std::string, ::std::function<bool(FileDescriptorProto*)>> DescMap;
    static const DescMap* google_protobuf = new DescMap{
        {
            "google/protobuf/any.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Any),
        },
        {
            "google/protobuf/api.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Api),
        },
        //        {
        //            "google/protobuf/compiler/plugin.proto",
        //            BUILD_FILE_DESCRIPTOR_PROTO(compiler::Version),
        //        },
        {
            "google/protobuf/descriptor.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(FileDescriptorProto),
        },
        {
            "google/protobuf/duration.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Duration),
        },
        {
            "google/protobuf/empty.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Empty),
        },
        {
            "google/protobuf/field_mask.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(FieldMask),
        },
        {
            "google/protobuf/source_context.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(SourceContext),
        },
        {
            "google/protobuf/struct.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Struct),
        },
        {
            "google/protobuf/timestamp.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Timestamp),
        },
        {
            "google/protobuf/type.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(Field),
        },
        {
            "google/protobuf/wrappers.proto",
            BUILD_FILE_DESCRIPTOR_PROTO(BoolValue),
        }};
    auto result = google_protobuf->find(filename);
    if (result == google_protobuf->end()) {
        return false;
    }
    return result->second(output);
}

bool loadFileDescriptorSet(FileDescriptorSet* files, const FileDescriptor* main, Errors* errors, std::unordered_set<std::string>* walk_set, PbParserOption* option) {
    typedef ::std::unordered_set<::std::string> std_string_set;
    static const std_string_set* google_protobuf = new std_string_set{
        "google/protobuf/any.proto",
        "google/protobuf/api.proto",
        "google/protobuf/compiler/plugin.proto",
        "google/protobuf/descriptor.proto",
        "google/protobuf/duration.proto",
        "google/protobuf/empty.proto",
        "google/protobuf/field_mask.proto",
        "google/protobuf/source_context.proto",
        "google/protobuf/struct.proto",
        "google/protobuf/timestamp.proto",
        "google/protobuf/type.proto",
        "google/protobuf/wrappers.proto",
    };
    if (option != nullptr && !option->with_google_protobuf_) {
        auto result = google_protobuf->find(main->name());
        if (result != google_protobuf->end()) {
            return true;
        }
    }
    auto exist = walk_set->find(main->name());
    if (exist != walk_set->end()) {
        return true;
    }
    if (main->name().find_first_of(""))
        walk_set->insert(main->name());
    auto newFileDesc = files->add_file();
    main->CopyTo(newFileDesc);
    if (option != nullptr) {
        if (option->with_source_code_info_) {
            main->CopySourceCodeInfoTo(newFileDesc);
        }
        if (option->with_json_tag_) {
            main->CopyJsonNameTo(newFileDesc);
        }
    }
    newFileDesc->set_syntax(main->SyntaxName(main->syntax()));
    for (int i = 0; i < main->dependency_count(); ++i) {
        auto dep = main->dependency(i);
        if (dep == nullptr) {
            errors->AddError("", -1, 0, "not found file");
            return false;
        }
        if (!loadFileDescriptorSet(files, dep, errors, walk_set, option)) {
            return false;
        }
    }
    return true;
}

#ifdef __cplusplus /* If this is a C++ compiler, use C linkage */
extern "C" {
#include "cgo.h"
#endif

struct ParserResult* ParsePBFile_C(const char* file, int size, struct ParserOption* option) {
    WrappedErrorCollector collector;
    PbParserOption cpp_option;
    cpp_option.message_type_ = PbParserOption::MessageType(option->message_type);
    cpp_option.require_syntax_identifier_ = bool(option->require_syntax_identifier);
    cpp_option.with_source_code_info_ = bool(option->with_source_code_info);
    cpp_option.with_json_tag_ = bool(option->with_json_tag);
    cpp_option.with_google_protobuf_ = bool(option->with_google_protobuf);

    auto desc = ParsePBFile(file, size_t(size), &collector, &cpp_option);
    auto result = new ParserResult{};

    result->desc_pointer = desc.release();
    result->errors_size = int(collector.errors.size());
    result->errors = new ParserErrorMessage[collector.errors.size()];

    for (int x = 0; x < int(collector.errors.size()); x++) {
        auto item = collector.errors.at(x);
        result->errors[x] = ParserErrorMessage{
            .line = item.line,
            .column = item.column,
            .message = parser::pb::dupString(item.message),
            .message_size = int(item.message.length()),
        };
    }
    return result;
}

struct ParserResult* ParseMultiPBFile_C(struct ParerFiles* files, struct ParserOption* option) {
    StringMap tree;
    tree.reserve(files->files_size);
    for (int x = 0; x < files->files_size; x++) {
        auto file = files->files[x];
        tree[file.filename] = {
            .data_ = file.content,
            .size_ = size_t(file.content_size),
        };
    }
    //    std::cout << "main: " << files->main_filename << std::endl;
    //    for (const auto& item : tree) {
    //        std::cout << "filename: " << item.first << std::endl;
    //        std::cout << "content: " << item.second.data_ << std::endl;
    //    }
    WrappedErrorCollector collector;
    PbParserOption cpp_option;
    for (int x = 0; x < option->include_name_size; x++) {
        cpp_option.include_path_.push_back(option->include_name[x]);
    }
    cpp_option.message_type_ = PbParserOption::MessageType(option->message_type);
    cpp_option.require_syntax_identifier_ = bool(option->require_syntax_identifier);
    cpp_option.with_source_code_info_ = bool(option->with_source_code_info);
    cpp_option.with_json_tag_ = bool(option->with_json_tag);
    cpp_option.with_google_protobuf_ = bool(option->with_google_protobuf);
    auto desc = ParseMultiPBFile(std::string(files->main_filename), tree, &collector, &cpp_option);
    auto result = new ParserResult{};
    result->desc_pointer = desc.release();
    if (collector.errors.size() > 0) {
        result->errors_size = int(collector.errors.size());
        result->errors = new ParserErrorMessage[collector.errors.size()];
    }
    for (int x = 0; x < int(collector.errors.size()); x++) {
        auto& item = collector.errors.at(x);
        result->errors[x] = ParserErrorMessage{
            .line = item.line,
            .column = item.column,
            .message = parser::pb::dupString(item.message),
            .message_size = int(item.message.length()),
            .filename = parser::pb::dupString(item.filename),
        };
    }
    return result;
}

void Delete_ParserResult(struct ParserResult* r) {
    for (int i = 0; i < r->errors_size; ++i) {
        delete r->errors[i].message;
        delete r->errors[i].filename;
    }
    delete[] r->errors;
    delete (std::string*)r->desc_pointer;
    delete r;
}

const char* ParserResult_GetDesc(struct ParserResult* r) {
    if (r->desc_pointer == nullptr) {
        return nullptr;
    }
    std::string* desc = (std::string*)r->desc_pointer;
    return desc->c_str();
}

int ParserResult_GetDescSize(struct ParserResult* r) {
    if (r->desc_pointer == nullptr) {
        return 0;
    }
    std::string* desc = (std::string*)r->desc_pointer;
    return int(desc->size());
}

#ifdef __cplusplus /* If this is a C++ compiler, end C linkage */
}
#endif

}  // namespace pb
}  // namespace parser
