//
// Created by ashy5000 on 6/5/24.
//

#include "BlockasmGenerator.h"

#include <iostream>
#include <sstream>

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
            Token intToken = tokens[i + 3];
            if(intToken.type != TokenType::int_lit) {
                std::cerr << "Expected integer literal." << std::endl;
                exit(EXIT_FAILURE);
            }
            Token closeParen = tokens[i + 4];
            if(closeParen.type != TokenType::close_paren) {
                std::cerr << "Expected ')'." << std::endl;
            }
            blockasm << "InitBfr 0x" << std::hex << nextAllocatedLocation << " 0x00000000" << std::endl;
            blockasm << "Exit 0x" << std::hex << nextAllocatedLocation << std::endl;
            nextAllocatedLocation++;
        }
    }
    return blockasm.str();
}
