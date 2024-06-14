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
#include "ParamsParser.h"
#include "Signature.h"
#include "SystemFunctions.h"
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
    for(const SystemFunction& func : SYSTEM_FUNCTIONS) {
        if(func.module == module && func.name == function) {
            std::string funcBlockasm = func.generateBlockasm(params, nextAllocatedLocation, vars);
            blockasm << funcBlockasm;
            return std::make_tuple(vars, 6);
        }
    }
    std::cerr << "Unknown module." << std::endl;
    exit(EXIT_FAILURE);
}
