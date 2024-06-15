//
// Created by ashy5000 on 6/13/24.
//

#include "ParamsParser.h"

#include <iostream>
#include <ostream>
#include <sstream>
#include <tuple>
#include <utility>

#include "ExpressionBlockasmGenerator.h"
#include "Signature.h"
#include "Type.h"

std::tuple<std::string, std::vector<int>, bool> ParamsParser::ParseParamsWithSignature(int &nextAllocatedLocation, const std::vector<Variable>& vars, const Signature& sig) {
    std::stringstream blockasm;
    std::vector<int> locations;
    std::vector<Type> types;
    for(Token param : params) {
        std::tuple<std::string, int, Type> expressionGenerationResult = ExpressionBlockasmGenerator::GenerateBlockasmFromExpression(params[0], nextAllocatedLocation, vars);
        std::string blockasmSection = std::get<0>(expressionGenerationResult);
        blockasm << blockasmSection;
        int location = std::get<1>(expressionGenerationResult);
        if(location >= nextAllocatedLocation) {
            nextAllocatedLocation = location + 1;
        }
        locations.emplace_back(location);
        Type type = std::get<2>(expressionGenerationResult);
        types.emplace_back(type);
    }
    if(!sig.CheckSignature(types)) {
        return std::make_tuple("", locations, false);
    }
    return std::make_tuple(blockasm.str(), locations, true);
}

std::tuple<std::string, std::vector<int>, Signature> ParamsParser::ParseParams(
    int &nextAllocatedLocation, const std::vector<Variable> &vars) {
    for(const Signature& sig : signatures) {
        std::tuple parseRes = ParseParamsWithSignature(nextAllocatedLocation, vars, sig);
        if(bool sigSucceeded = std::get<2>(parseRes); !sigSucceeded) {
            continue;
        }
        std::string blockasm = std::get<0>(parseRes);
        std::vector locations = std::get<1>(parseRes);
        return std::make_tuple(blockasm, locations, sig);
    }
    std::cerr << "No matching signature for function" << std::endl;
    exit(EXIT_FAILURE);
}

ParamsParser::ParamsParser(std::vector<Token> params_p, std::vector<Signature> signatures_p) : params(std::move(params_p)), signatures(std::move(signatures_p)) {}
