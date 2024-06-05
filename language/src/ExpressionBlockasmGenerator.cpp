//
// Created by ashy5000 on 6/5/24.
//

#include "ExpressionBlockasmGenerator.h"

std::string ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(Token expression, int nextAllocatedLocation) {
    if(expression.type != TokenType::expr) {
        std::cerr << "Expected expression when generating Blockasm." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children.size() == 0) {
        std::cerr << "Empty expression not allowed." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children.size() == 1) {
        if(expression.children[0].type == TokenType::int_lit) {
            std::stringstream blockasm;
            blockasm << "InitBfr 0x" << std::hex << nextAllocatedLocation << " " << expression.children[0].value;
            return blockasm.str();
        } else {
            std::cerr << "Unknown expression." << std::endl;
        }
    }
    std::cerr << "Unknown expression." << std::endl;
    exit(EXIT_FAILURE);
}
