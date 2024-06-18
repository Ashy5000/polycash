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
#include "Linker.h"
#include "ParamsParser.h"
#include "Signature.h"
#include "SystemFunctions.h"
#include "Variable.h"

BlockasmGenerator::BlockasmGenerator(std::vector<Token> tokens_p) {
    tokens = std::move(tokens_p);
    blockasm = {};
    blockasm << ";^^^^BEGIN_SOURCE^^^^" << std::endl;
}


std::string BlockasmGenerator::GenerateBlockasm() {
    std::vector<Variable> vars;
    int nextAllocatedLocation = 0x00001000;
    auto l = Linker({"string.blockasm"});
    for(int i = 0; i < tokens.size(); i++) {
        if(const Token token = tokens[i]; token.type == TokenType::system_at) {
            std::tuple tuple = GenerateSystemFunctionBlockasm(i, nextAllocatedLocation, vars, l);
            std::vector<Variable> newVars = std::get<0>(tuple);
            vars.insert(vars.end(), newVars.begin(), newVars.end());
            const int tokensConsumed = std::get<1>(tuple);
            i += tokensConsumed;
        } else if(token.type == TokenType::identifier) {
            if(tokens[i + 1].type == TokenType::eq) {
                if(tokens[i + 2].type == TokenType::eq) {
                    // e.g. newVar == 3
                    std::string varName = token.value;
                    std::tuple exprTuple = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(tokens[i + 4], nextAllocatedLocation, vars);
                    std::string exprBlockasm = std::get<0>(exprTuple);
                    blockasm << exprBlockasm;
                    int exprLoc = std::get<1>(exprTuple);
                    if(exprLoc >= nextAllocatedLocation) {
                        nextAllocatedLocation = exprLoc + 1;
                    }
                    Type type = std::get<2>(exprTuple);
                    Variable var = Variable(varName, exprLoc, type);
                    vars.emplace_back(var);
                    i += 6;
                } else {
                    // e.g. existingVar = 5
                    std::string varName = token.value;
                    std::tuple exprTuple = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(tokens[i + 3], nextAllocatedLocation, vars);
                    std::string exprBlockasm = std::get<0>(exprTuple);
                    blockasm << exprBlockasm;
                    int exprLoc = std::get<1>(exprTuple);
                    if(exprLoc >= nextAllocatedLocation) {
                        nextAllocatedLocation = exprLoc + 1;
                    }
                    Type type = std::get<2>(exprTuple);
                    for(const Variable &var : vars) {
                        if(var.name == varName) {
                            if(var.type != type) {
                                std::cerr << "Expression type does not match variable type." << std::endl;
                                exit(EXIT_FAILURE);
                            }
                            blockasm << "CpyBfr 0x" << std::setfill('0') << std::setw(8) << exprLoc << " 0x";
                            blockasm << std::setfill('0') << std::setw(8) << var.location << " 0x00000000" << std::endl;
                            break;
                        }
                    }
                    i += 5;
                }
            }
        } else if(token.type == TokenType::div) {
            if(tokens[i + 1].type == TokenType::identifier) {
                l.InjectIfNotPresent(tokens[i + 1].value, blockasm);
            }
        } else if(token.type == TokenType::excl) {
            if(tokens[i + 1].type == TokenType::identifier) {
                std::string functionName = tokens[i + 1].value;
                if(tokens[i + 2].type != TokenType::open_paren) {
                    std::cerr << "Expected '('" << std::endl;
                    exit(EXIT_FAILURE);
                }
                std::vector<Token> params;
                for(int j = i + 3; j < tokens.size(); j++) {
                    if(tokens[j].type == TokenType::expr) {
                        params.emplace_back(tokens[j]);
                    }
                    if(tokens[j].type == TokenType::newline) {
                        break;
                    }
                }
                std::vector<int> paramLocs;
                for(Token param : params) {
                    std::tuple exprTuple = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(param, nextAllocatedLocation, vars);
                    std::string exprBlockasm = std::get<0>(exprTuple);
                    blockasm << exprBlockasm;
                    int exprLoc = std::get<1>(exprTuple);
                    if(exprLoc >= nextAllocatedLocation) {
                        nextAllocatedLocation = exprLoc + 1;
                    }
                    paramLocs.emplace_back(exprLoc);
                }
                std::string functionCallBlockasm = l.CallFunction(functionName, paramLocs);
                blockasm << functionCallBlockasm;
            }
        }
    }
    Linker::SkipLibs(blockasm);
    std::string blockasmStr = blockasm.str();
    return blockasmStr;
}

std::tuple<std::vector<Variable>, int> BlockasmGenerator::GenerateSystemFunctionBlockasm(const int i, int &nextAllocatedLocation, std::vector<Variable> vars, Linker l) {
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
            std::string funcBlockasm = func.generateBlockasm(params, nextAllocatedLocation, vars, l);
            blockasm << funcBlockasm;
            return std::make_tuple(vars, 6);
        }
    }
    std::cerr << "Unknown module." << std::endl;
    exit(EXIT_FAILURE);
}
