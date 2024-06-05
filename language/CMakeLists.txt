cmake_minimum_required(VERSION 3.20)

project(polycash_lang)

set(CMAKE_CXX_STANDARD 20)

add_executable(plc src/main.cpp
        src/main.cpp
        src/TokenType.h
        src/Token.h
        src/Parser.cpp
        src/Parser.h
        src/BlockasmGenerator.cpp
        src/BlockasmGenerator.h
)