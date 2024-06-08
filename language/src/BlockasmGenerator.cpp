//
// Created by ashy5000 on 6/5/24.
//

#include "BlockasmGenerator.h"

#include <iomanip>
#include <iostream>
#include <sstream>
#include <tuple>
#include <utility>

#include "ExpressionBlockasmGenerator.h"
#include "Variable.h"

BlockasmGenerator::BlockasmGenerator(std::vector<Token> tokens_p) {
    tokens = std::move(tokens_p);
    blockasm = {};
}


std::string BlockasmGenerator::GenerateBlockasm() {
    std::vector<Variable> vars;
    int nextAllocatedLocation = 1;
    for(int i = 0; i < tokens.size(); i++) {
        if(Token token = tokens[i]; token.type == TokenType::system_at) {
            std::vector<Variable> newVars = GenerateSystemFunctionBlockasm(i, nextAllocatedLocation, vars);
            vars.insert(vars.end(), newVars.begin(), newVars.end());
            i += 6;
        }
    }
    return blockasm.str();
}

std::vector<Variable> BlockasmGenerator::GenerateSystemFunctionBlockasm(int i, int &nextAllocatedLocation, std::vector<Variable> vars) {
    Token identifier = tokens[i + 1];
    if(identifier.type != TokenType::identifier) {
        std::cerr << "System at (@) must be followed by an identifier." << std::endl;
        exit(EXIT_FAILURE);
    }
    Token openParen = tokens[i + 2];
    if(openParen.type != TokenType::open_paren) {
        std::cerr << "System call identifier must be followed by '('." << std::endl;
        exit(EXIT_FAILURE);
    }
    Token exprToken = tokens[i + 3];
    if(exprToken.type != TokenType::expr) {
        std::cerr << "Expected expression." << std::endl;
        exit(EXIT_FAILURE);
    }
    Token closeParen = tokens[i + 4];
    if(closeParen.type != TokenType::close_paren) {
        std::cerr << "Expected ')'." << std::endl;
        exit(EXIT_FAILURE);
    }
    std::string delimiter = "::";
    int delimiterPos = identifier.value.find(delimiter);
    if(delimiterPos == identifier.value.npos) {
        std::cerr << "Invalid system function format." << std::endl;
        exit(EXIT_FAILURE);
    }
    std::string module = identifier.value.substr(0, delimiterPos);
    std::string function = identifier.value.substr(delimiterPos + 2);
    if(module == "contract") {
        if(function == "exit") {
            ExpressionBlockasmGenerator generator;
            std::tuple<std::string, int> expressionGenerationResult = generator.GenerateBlockasmFromExpression(exprToken, nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
            int location = std::get<1>(expressionGenerationResult);
            blockasm << expressionBlockasm << std::endl;
            blockasm << "Exit 0x" << std::setfill('0') << std::setw(8) << std::hex << location << std::endl;
            if(location >= nextAllocatedLocation) {
                nextAllocatedLocation = location + 1;
            }
            Token semiToken = tokens[i + 5];
            if(semiToken.type != TokenType::semi) {
                std::cerr << "Expected semicolon." << std::endl;
                exit(EXIT_FAILURE);
            }
            Token newlineToken = tokens[i + 6];
            if(newlineToken.type != TokenType::newline) {
                std::cerr << "Unexpected token after semicolon." << std::endl;
                exit(EXIT_FAILURE);
            }
        } else {
            std::cerr << "Unknown system function " << identifier.value << "." << std::endl;
            exit(EXIT_FAILURE);
        }
    } else if(module == "memory") {
        if(function == "alloc") {
            Variable var = Variable(exprToken.children[0].value, nextAllocatedLocation, Type::type_placeholder);
            vars.emplace_back(var);
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation++ << " 0x00000000" << std::endl;
        } else if(function == "free") {
            int indexToRemove = -1;
            for(int j = 0; j < vars.size(); j++) {
                Variable var = vars[j];
                if(var.name == exprToken.children[0].value) {
                    indexToRemove = j;
                    break;
                }
            }
            if(indexToRemove == -1) {
                std::cerr << "Cannot free undefined variable." << std::endl;
                exit(EXIT_FAILURE);
            }
            blockasm << "Free 0x" << std::setfill('0') << std::setw(8) << std::hex << vars[indexToRemove].location << " 0x00000000" << std::endl;
            vars.erase(vars.begin() + indexToRemove);
        } else {
            std::cerr << "Unknown system function " << function << "." << std::endl;
            exit(EXIT_FAILURE);
        }
    } else if(module == "io") {
        if(function == "printf") {
            ExpressionBlockasmGenerator generator;
            std::tuple<std::string, int> expressionGenerationResult = generator.GenerateBlockasmFromExpression(exprToken, nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
            int expressionLocation = std::get<1>(expressionGenerationResult);
            blockasm << expressionBlockasm << std::endl;
            blockasm << "Stdout 0x" << std::setfill('0') << std::setw(8) << std::hex << expressionLocation << " 0x00000000" << std::endl;
            blockasm << "Free 0x" << std::setfill('0') << std::setw(8) << std::hex << expressionLocation << " 0x00000000" << std::endl;
        }
    } else {
        std::cerr << "Unknown module." << std::endl;
        exit(EXIT_FAILURE);
    }
    return vars;
}
