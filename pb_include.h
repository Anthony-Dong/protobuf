//
// Created by bytedance on 2023/5/11.
//

#ifndef PROTOBUF_PB_INCLUDE_H
#define PROTOBUF_PB_INCLUDE_H
#include <google/protobuf/any.pb.h>
#include <google/protobuf/api.pb.h>
#include <google/protobuf/descriptor.pb.h>
#include <google/protobuf/duration.pb.h>
#include <google/protobuf/empty.pb.h>
#include <google/protobuf/field_mask.pb.h>
#include <google/protobuf/source_context.pb.h>
#include <google/protobuf/struct.pb.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/type.pb.h>
#include <google/protobuf/wrappers.pb.h>

#include <functional>
#include <string>
#include <unordered_map>

namespace parser {
namespace pb {

#undef BUILD_FILE_DESCRIPTOR_PROTO
#define BUILD_FILE_DESCRIPTOR_PROTO(TypeName)           \
    [](google::protobuf::FileDescriptorProto* output) { \
        TypeName::descriptor()->file()->CopyTo(output); \
        return true;                                    \
    }

}  // namespace pb
}  // namespace parser

#endif  // PROTOBUF_PB_INCLUDE_H
