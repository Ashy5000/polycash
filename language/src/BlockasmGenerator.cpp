//
// Created by ashy5000 on 6/5/24.
//

#include "BlockasmGenerator.h"

#include <iostream>
#include <sstream>
#include <tuple>

#include "ExpressionBlockasmGenerator.h"
#include "Variable.h"

std::string BlockasmGenerator::GenerateBlockasm(std::vector<Token> tokens) {
    std::stringstream blockasm;
    std::vector<Variable> vars;
    int nextAllocatedLocation = 0;
    for(int i = 0; i < tokens.size(); i++) {
        Token token = tokens[i];
        if(token.type == TokenType::system_at) {
            std::tuple<std::string, std::vector<Variable>> generatedTuple = GenerateSystemFunctionBlockasm(tokens, i, nextAllocatedLocation);
            blockasm << std::get<0>(generatedTuple);
            std::vector<Variable> newVars = std::get<1>(generatedTuple);
            vars.insert(vars.end(), newVars.begin(), newVars.end());
            i += 6;
        }
    }
    return blockasm.str();
}

std::tuple<std::string, std::vector<Variable>> BlockasmGenerator::GenerateSystemFunctionBlockasm(std::vector<Token> tokens, int i, int &nextAllocatedLocation) {
    std::vector<Variable> vars;
    std::stringstream blockasm;
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
    }
    if(identifier.value == "exit") {
        ExpressionBlockasmGenerator generator;
        std::tuple<std::string, int> expressionGenerationResult = generator.GenerateBlockasmFromExpression(exprToken, nextAllocatedLocation, vars);
        std::string expressionBlockasm = std::get<0>(expressionGenerationResult);
        int location = std::get<1>(expressionGenerationResult);
        blockasm << expressionBlockasm << std::endl;
        blockasm << "Exit 0x" << std::hex << location << std::endl;
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
    } else if(identifier.value == "alloc") {
        Variable var = Variable(exprToken.children[0].value, nextAllocatedLocation);
        vars.emplace_back(var);
        blockasm << "InitBfr 0x" << std::hex << nextAllocatedLocation++ << std::endl;
    } else {
        std::cerr << "Unknown system function." << std::endl;
        exit(EXIT_FAILURE);
    }
    return std::make_tuple(blockasm.str(), vars);
}
