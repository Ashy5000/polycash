//
// Created by ashy5000 on 6/5/24.
//

#include "ExpressionBlockasmGenerator.h"

#include <iomanip>
#include <iostream>
#include <sstream>
#include <tuple>

#include "Operator.h"
#include "OperatorType.h"
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
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " " << expression.children[0].value;
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
        exit(EXIT_FAILURE);
    }
    OperatorType type = OperatorType::type_placeholder;
    int i = 0;
    for(;i < expression.children.size(); i++) {
        Token t = expression.children[i];
        if(t.type == TokenType::concat) {
            type = OperatorType::concat;
            break;
        }
    }
    std::vector preOperatorTokens(expression.children.begin(), expression.children.begin() + i);
    std::vector postOperatorTokens(expression.children.begin() + i + 1, expression.children.end());
    Token preOperatorExpr = Token({TokenType::expr, {}});
    preOperatorExpr.children = preOperatorTokens;
    Token postOperatorExpr = Token({TokenType::expr, {}});
    postOperatorExpr.children = postOperatorTokens;
    if(type == OperatorType::concat) {
        std::stringstream blockasm;
        std::tuple exprATuple = GenerateBlockasmFromExpression(preOperatorExpr, nextAllocatedLocation + 1, vars);
        std::string exprABlockasm = std::get<0>(exprATuple);
        blockasm << exprABlockasm;
        int exprALoc = std::get<1>(exprATuple);
        std::tuple exprBTuple = GenerateBlockasmFromExpression(postOperatorExpr, nextAllocatedLocation + 1, vars);
        std::string exprBBlockasm = std::get<0>(exprATuple);
        blockasm << exprBBlockasm;
        int exprBLoc = std::get<1>(exprATuple);
        blockasm << "App 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x" << exprBLoc << " 0x" <<  nextAllocatedLocation + 1 << " 0x00000000" << std::endl;
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1);
    }
    std::cerr << "Unknown expression." << std::endl;
    exit(EXIT_FAILURE);
}
