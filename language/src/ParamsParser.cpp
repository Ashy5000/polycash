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

std::tuple<std::string, std::vector<int>> ParamsParser::ParseParams(int &nextAllocatedLocation, const std::vector<Variable>& vars) {
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
        std::cerr << "Incorrect signature for function." << std::endl;
        exit(EXIT_FAILURE);
    }
    return std::make_tuple(blockasm.str(), locations);
}

ParamsParser::ParamsParser(std::vector<Token> params_p, Signature sig_p) : params(std::move(params_p)), sig(std::move(sig_p)) {
}
