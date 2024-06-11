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
        if(const Token token = tokens[i]; token.type == TokenType::system_at) {
            std::tuple tuple = GenerateSystemFunctionBlockasm(i, nextAllocatedLocation, vars);
            std::vector<Variable> newVars = std::get<0>(tuple);
            vars.insert(vars.end(), newVars.begin(), newVars.end());
            const int tokensConsumed = std::get<1>(tuple);
            i += tokensConsumed;
        }
    }
    return blockasm.str();
}

std::tuple<std::vector<Variable>, int> BlockasmGenerator::GenerateSystemFunctionBlockasm(const int i, int &nextAllocatedLocation, std::vector<Variable> vars) {
    Token identifier = tokens[i + 1];
    if(identifier.type != TokenType::identifier) {
        std::cerr << "System at (@) must be followed by an identifier." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(Token openParen = tokens[i + 2]; openParen.type != TokenType::open_paren) {
        std::cerr << "System call identifier must be followed by '('." << std::endl;
        exit(EXIT_FAILURE);
    }
    std::vector<Token> params;
    std::vector<Token> currentExprTokens;
    for(int j = 0; j < tokens[i + 3].children.size(); j++) {
        Token t = tokens[i + 3].children[j];
        if(t.type == TokenType::comma) {
            auto expr = Token(TokenType::expr, {});
            for(const Token& exprT : currentExprTokens) {
                expr.children.emplace_back(exprT);
            }
            params.emplace_back(expr);
            currentExprTokens.clear();
            continue;
        }
        currentExprTokens.emplace_back(t);
        if(j == tokens[i + 3].children.size() - 1) {
            auto expr = Token(TokenType::expr, {});
            for(const Token& exprT : currentExprTokens) {
                expr.children.emplace_back(exprT);
            }
            params.emplace_back(expr);
            currentExprTokens.clear();
            break;
        }
    }
    if(Token semiToken = tokens[i + 5]; semiToken.type != TokenType::semi) {
        std::cerr << "Expected semicolon." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(Token newlineToken = tokens[i + 6]; newlineToken.type != TokenType::newline) {
        std::cerr << "Unexpected token after semicolon." << std::endl;
        exit(EXIT_FAILURE);
    }
    std::string delimiter = "::";
    auto delimiterPos = identifier.value.find(delimiter);
    if(delimiterPos == std::string::npos) {
        std::cerr << "Invalid system function format." << std::endl;
        exit(EXIT_FAILURE);
    }
    std::string module = identifier.value.substr(0, delimiterPos);
    std::string function = identifier.value.substr(delimiterPos + 2);
    if(module == "contract") {
        if(function == "exit") {
            std::tuple<std::string, int> expressionGenerationResult = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(params[0], nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
            int location = std::get<1>(expressionGenerationResult);
            blockasm << expressionBlockasm << std::endl;
            blockasm << "ExitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << location << std::endl;
            if(location >= nextAllocatedLocation) {
                nextAllocatedLocation = location + 1;
            }
        } else {
            std::cerr << "Unknown system function " << identifier.value << "." << std::endl;
            exit(EXIT_FAILURE);
        }
    } else if(module == "memory") {
        if(function == "alloc") {
            auto var = Variable(params[0].children[0].value, nextAllocatedLocation, Type::type_placeholder);
            vars.emplace_back(var);
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation++ << " 0x00000000" << std::endl;
        } else if(function == "free") {
            int indexToRemove = -1;
            for(int j = 0; j < vars.size(); j++) {
                if(Variable var = vars[j]; var.name == params[0].children[0].value) {
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
        } else if(function == "set") {
            int indexToRename = -1;
            for(int j = 0; j < vars.size(); j++) {
                if(Variable var = vars[j]; var.name == params[0].children[0].value) {
                    indexToRename = j;
                    break;
                }
            }
            if(indexToRename == -1) {
                std::cerr << "Cannot set undefined variable." << std::endl;
                exit(EXIT_FAILURE);
            }
            char* end;
            int val = std::strtol(params[1].children[0].value.c_str(), end, 10);
            if(errno == ERANGE) {
                std::cerr << "Expected integer as value" << std::endl;
            }
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << vars[indexToRename].location << " 0x";
            blockasm << std::setfill('0') << std::setw(16) << std::hex << val << " 0x00000000" << std::endl;
        } else {
            std::cerr << "Unknown system function " << function << "." << std::endl;
            exit(EXIT_FAILURE);
        }
    } else if(module == "io") {
        if(function == "print") {
            std::tuple<std::string, int> expressionGenerationResult = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(params[0], nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
            int expressionLocation = std::get<1>(expressionGenerationResult);
            if(expressionLocation >= nextAllocatedLocation) {
                nextAllocatedLocation = expressionLocation + 1;
            }
            blockasm << expressionBlockasm << std::endl;
            blockasm << "Stdout 0x" << std::setfill('0') << std::setw(8) << std::hex << expressionLocation << " 0x00000000" << std::endl;
        } else if(function == "err") {
            std::tuple<std::string, int> expressionGenerationResult = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(params[0], nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
            int expressionLocation = std::get<1>(expressionGenerationResult);
            if(expressionLocation >= nextAllocatedLocation) {
                nextAllocatedLocation = expressionLocation + 1;
            }
            blockasm << expressionBlockasm << std::endl;
            blockasm << "Stderr 0x" << std::setfill('0') << std::setw(8) << std::hex << expressionLocation << " 0x00000000" << std::endl;
        }
    } else {
        std::cerr << "Unknown module." << std::endl;
        exit(EXIT_FAILURE);
    }
    return std::make_tuple(vars, 6);
}
