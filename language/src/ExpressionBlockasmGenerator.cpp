//
// Created by ashy5000 on 6/5/24.
//

#include "ExpressionBlockasmGenerator.h"

#include <iomanip>
#include <iostream>
#include <iterator>
#include <sstream>
#include <tuple>

#include "OperatorType.h"
#include "Variable.h"

std::tuple<std::string, int, Type> ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(const Token &expression, int nextAllocatedLocation, const std::vector<Variable>& vars) {
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
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x00000000" << std::endl;
            int val = stoi(expression.children[0].value);
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x";
            blockasm << std::setfill('0') << std::setw(16) << std::hex << val << " 0x00000000" << std::endl;
            return std::make_tuple(blockasm.str(), nextAllocatedLocation, Type::uint64);
        }
        if(expression.children[0].type == TokenType::string_lit) {
            std::stringstream blockasm;
            std::stringstream buffer(expression.children[0].value);
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x00000000" << std::endl;
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x";
            std::istreambuf_iterator it( buffer.rdbuf( ) );
            std::istreambuf_iterator<char> end; // eof
            std::stringstream out;
            out << std::hex;
            std::copy(it, end, std::ostream_iterator<int>(out));
            blockasm << out.str();
            blockasm << " 0x00000000" << std::endl;
            return std::make_tuple(blockasm.str(), nextAllocatedLocation, Type::string);
        }
        if(expression.children[0].type == TokenType::identifier) {
            auto referencedVar = Variable("", 0, Type::type_placeholder);
            for(const Variable& var : vars) {
                if(var.name == expression.children[0].value) {
                    referencedVar = var;
                }
            }
            return std::make_tuple("", referencedVar.location, referencedVar.type);
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
        if(t.type == TokenType::add) {
            type = OperatorType::add;
            break;
        }
        if(t.type == TokenType::sub) {
            type = OperatorType::sub;
            break;
        }
        if(t.type == TokenType::mul) {
            type = OperatorType::mul;
            break;
        }
        if(t.type == TokenType::div) {
            type = OperatorType::div;
            break;
        }
    }
    std::vector preOperatorTokens(expression.children.begin(), expression.children.begin() + i);
    std::vector postOperatorTokens(expression.children.begin() + i + 1, expression.children.end());
    Token preOperatorExpr = Token({TokenType::expr, {}});
    preOperatorExpr.children = preOperatorTokens;
    Token postOperatorExpr = Token({TokenType::expr, {}});
    postOperatorExpr.children = postOperatorTokens;
    std::stringstream blockasm;
    std::tuple exprATuple = GenerateBlockasmFromExpression(preOperatorExpr, nextAllocatedLocation, vars);
    std::string exprABlockasm = std::get<0>(exprATuple);
    blockasm << exprABlockasm;
    int exprALoc = std::get<1>(exprATuple);
    if(exprALoc >= nextAllocatedLocation) {
        nextAllocatedLocation = exprALoc + 1;
    }
    std::tuple exprBTuple = GenerateBlockasmFromExpression(postOperatorExpr, nextAllocatedLocation, vars);
    std::string exprBBlockasm = std::get<0>(exprBTuple);
    blockasm << exprBBlockasm;
    int exprBLoc = std::get<1>(exprBTuple);
    if(exprALoc >= nextAllocatedLocation) {
        nextAllocatedLocation = exprALoc + 1;
    }
    blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000" << std::endl;
    if(type == OperatorType::concat) {
        blockasm << "App 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::add) {
        blockasm << "Add 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::sub) {
        blockasm << "Sub 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::mul) {
        blockasm << "Mul 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::div) {
        blockasm << "Div 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(blockasm.str(), nextAllocatedLocation + 1, Type::uint64);
    }
    std::cerr << "Unknown expression." << std::endl;
    exit(EXIT_FAILURE);
}
