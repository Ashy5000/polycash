//
// Created by ashy5000 on 6/5/24.
//

#include "ExpressionBlockasmGenerator.h"

#include <iostream>
#include <sstream>
#include <tuple>

#include "Variable.h"

std::tuple<std::string, int> ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(const Token &expression, const int nextAllocatedLocation, const std::vector<Variable>& vars) {
    if(expression.type != TokenType::expr) {
        std::cerr << "Expected expression when generating Blockasm." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children.empty()) {
        std::cerr << "Empty expression not allowed." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children.size() == 1) {
        if(expression.children[0].type == TokenType::int_lit) {
            std::stringstream blockasm;
            blockasm << "InitBfr 0x" << std::hex << nextAllocatedLocation << " " << expression.children[0].value;
            return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1);
        }
        if(expression.children[0].type == TokenType::identifier) {
            auto referencedVar = Variable("", 0, Type::type_placeholder);
            for(const Variable& var : vars) {
                if(var.name == expression.children[0].value) {
                    referencedVar = var;
                }
            }
            return std::make_tuple("", referencedVar.location);
        }
        std::cerr << "Unknown expression." << std::endl;
    }
    std::cerr << "Unknown expression." << std::endl;
    exit(EXIT_FAILURE);
}
