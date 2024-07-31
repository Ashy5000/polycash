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

std::tuple<int, Type> ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(Token expression, int nextAllocatedLocation, std::vector<Variable>& vars, std::stringstream &blockasm, Linker &l) {
    if(expression.type != TokenType::expr) {
        std::cerr << "Expected expression when generating Blockasm." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children.empty()) {
        std::cerr << "Empty expression not allowed." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children[0].type == TokenType::open_paren) {
        auto first = expression.children.begin() + 1;
        auto last = expression.children.end() - 1;
        expression.children = std::vector(first, last);
    }
    if(expression.children.size() == 1) {
        if (expression.children[0].type == TokenType::int_lit) {
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation
                     << " 0x00000000" << std::endl;
            int val = stoi(expression.children[0].value);
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x";
            blockasm << std::setfill('0') << std::setw(16) << std::hex << val << " 0x00000000" << std::endl;
            return std::make_tuple(nextAllocatedLocation, Type::uint64);
        }
        if (expression.children[0].type == TokenType::string_lit) {
            std::stringstream buffer(expression.children[0].value);
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation
                     << " 0x00000000" << std::endl;
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation << " 0x";
            std::istreambuf_iterator it(buffer.rdbuf());
            std::istreambuf_iterator<char> end; // eof
            std::stringstream out;
            out << std::hex;
            std::copy(it, end, std::ostream_iterator<int>(out));
            blockasm << out.str();
            blockasm << " 0x00000000" << std::endl;
            return std::make_tuple(nextAllocatedLocation, Type::string);
        }
        if (expression.children[0].type == TokenType::identifier) {
            auto referencedVar = Variable("", 0, Type::type_placeholder);
            for (const Variable &var: vars) {
                if (var.name == expression.children[0].value) {
                    referencedVar = var;
                    break;
                }
            }
            return std::make_tuple(referencedVar.location, referencedVar.type);
        }
        if(expression.children[0].type == TokenType::expr) {
            std::tuple exprTuple = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(expression.children[0], nextAllocatedLocation, vars, blockasm, l);
            return exprTuple;
        }
        std::cerr << "Unknown expression." << std::endl;
        exit(EXIT_FAILURE);
    }
    if(expression.children[0].type == TokenType::excl) {
        auto tokens = expression.children;
        if(tokens[1].type != TokenType::identifier) {
            std::cerr << "Expected identifier." << std::endl;
            exit(EXIT_FAILURE);
        }
        if(tokens[2].type != TokenType::open_paren) {
            std::cerr << "Expected '('" << std::endl;
            exit(EXIT_FAILURE);
        }
        std::vector<Token> params;
        std::string functionName = tokens[1].value;
        for(auto & token : tokens[3].children) {
            if(token.type == TokenType::expr) {
                params.emplace_back(token);
            }
            if(token.type == TokenType::newline) {
                break;
            }
        }
        std::vector<int> paramLocs;
        for(const Token& param : params) {
            std::tuple exprTuple = GenerateBlockasmFromExpression(param, nextAllocatedLocation, vars, blockasm, l);
            int exprLoc = std::get<0>(exprTuple);
            if(exprLoc >= nextAllocatedLocation) {
                nextAllocatedLocation = exprLoc + 1;
            }
            paramLocs.emplace_back(exprLoc);
        }
        std::tuple functionCallTuple = l.CallFunction(functionName, paramLocs, vars);
        std::string functionCallBlockasm = std::get<0>(functionCallTuple);
        blockasm << functionCallBlockasm;
        Type t = std::get<1>(functionCallTuple);
        return std::make_tuple(0x00000001, t);
    }
    OperatorType type = OperatorType::type_placeholder;
    int i = 0;
    for(;i < expression.children.size(); i++) {
        Token t = expression.children[i];
        type = OperatorTypeFromToken(t);
        if(type != OperatorType::type_placeholder) {
            break;
        }
    }
    std::vector preOperatorTokens(expression.children.begin(), expression.children.begin() + i);
    std::vector postOperatorTokens(expression.children.begin() + i + 1, expression.children.end());
    Token preOperatorExpr = Token({TokenType::expr, {}});
    preOperatorExpr.children = preOperatorTokens;
    Token postOperatorExpr = Token({TokenType::expr, {}});
    postOperatorExpr.children = postOperatorTokens;
    std::tuple exprATuple = GenerateBlockasmFromExpression(preOperatorExpr, nextAllocatedLocation, vars, blockasm, l);
    int exprALoc = std::get<0>(exprATuple);
    if(exprALoc >= nextAllocatedLocation) {
        nextAllocatedLocation = exprALoc + 1;
    }
    std::tuple exprBTuple = GenerateBlockasmFromExpression(postOperatorExpr, nextAllocatedLocation, vars, blockasm, l);
    int exprBLoc = std::get<0>(exprBTuple);
    if(exprALoc >= nextAllocatedLocation) {
        nextAllocatedLocation = exprALoc + 1;
    }
    blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000" << std::endl;
    if(type == OperatorType::concat) {
        blockasm << "App 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::add) {
        blockasm << "Add 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::sub) {
        blockasm << "Sub 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::mul) {
        blockasm << "Mul 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::div) {
        blockasm << "Div 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::uint64);
    }
    if(type == OperatorType::eq) {
        blockasm << "Eq 0x" << std::setfill('0') << std::setw(8) << std::hex << exprALoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << exprBLoc << " 0x";
        blockasm << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation + 1 << " 0x00000000";
        return std::make_tuple(nextAllocatedLocation + 1, Type::boolean);
    }
    std::cerr << "Unknown expression." << std::endl;
    exit(EXIT_FAILURE);
}
