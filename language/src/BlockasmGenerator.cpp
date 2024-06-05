//
// Created by ashy5000 on 6/5/24.
//

#include "BlockasmGenerator.h"

#include <iostream>
#include <sstream>

#include "ExpressionBlockasmGenerator.h"

std::string BlockasmGenerator::GenerateBlockasm(std::vector<Token> tokens) {
    std::stringstream blockasm;
    int nextAllocatedLocation = 0;
    for(int i = 0; i < tokens.size(); i++) {
        Token token = tokens[i];
        if(token.type == TokenType::system_at) {
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
            ExpressionBlockasmGenerator generator;
            std::string expressionBlockasm = generator.GenerateBlockasmFromExpression(exprToken, nextAllocatedLocation);
            blockasm << expressionBlockasm << std::endl;
            blockasm << "Exit 0x" << std::hex << nextAllocatedLocation << std::endl;
            nextAllocatedLocation++;
            Token semiToken = tokens[i + 5];
            if(semiToken.type != TokenType::semi) {
                std::cerr << "Expected semicolon." << std::endl;
            }
            Token newlineToken = tokens[i + 6];
            if(newlineToken.type != TokenType::newline) {
                std::cerr << "Unexpected token after semicolon." << std::endl;
            }
            i += 6;
        }
    }
    return blockasm.str();
}
