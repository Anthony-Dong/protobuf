cmake_minimum_required(VERSION 3.0.0)
set(CMAKE_CXX_STANDARD 11)

message(CMAKE_VERSION: ${CMAKE_VERSION})
message(CMAKE_SYSTEM: ${CMAKE_SYSTEM})

project(lib_parser)

add_library(${PROJECT_NAME} SHARED pb_parser.cpp pb_parser.h cgo.h)

set(protobuf_link_directory ${CMAKE_CURRENT_SOURCE_DIR}/deps/darwin_x86_64)
if(${CMAKE_SYSTEM} MATCHES "Linux")
    set(protobuf_link_directory ${CMAKE_CURRENT_SOURCE_DIR}/deps/linux_x86_64)
endif()

target_include_directories(${PROJECT_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/deps/include)
if (${CMAKE_VERSION} VERSION_LESS 3.13)
    link_directories(${protobuf_link_directory})
else ()
    target_link_directories(${PROJECT_NAME} PUBLIC ${protobuf_link_directory})
endif ()
target_link_libraries(${PROJECT_NAME} PUBLIC protobuf)
