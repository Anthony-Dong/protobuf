cmake_minimum_required(VERSION 3.0.0)

project(cgo_demo)

set(CMAKE_CXX_STANDARD 11)

set(PROTOBUF_DIR ${PROJECT_SOURCE_DIR}/../../)

add_subdirectory(${PROJECT_SOURCE_DIR}/../../ protobuf)

## cpp
add_executable(${PROJECT_NAME} main.cpp)
target_include_directories(${PROJECT_NAME} PRIVATE ${PROTOBUF_DIR})
target_link_libraries(${PROJECT_NAME} PUBLIC lib_parser)

## c
add_executable(cgo_c_demo main.c)
#target_include_directories(cgo_c_demo PRIVATE ${PROTOBUF_DIR})
#target_link_libraries(cgo_c_demo PUBLIC lib_parser)