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
        src/ExpressionBlockasmGenerator.cpp
        src/ExpressionBlockasmGenerator.h
        src/Variable.h
        src/Type.h
        src/Operator.h
        src/OperatorType.h
        src/Signature.cpp
        src/Signature.h
        src/ParamsParser.cpp
        src/ParamsParser.h
        src/SystemFunction.h
        src/SystemFunctions.h
        src/BlockasmLib.cpp
        src/BlockasmLib.h
        src/Function.h
        src/Linker.cpp
        src/Linker.h
        src/InjectedFunction.h
        src/LabelManager.cpp
        src/LabelManager.h
)