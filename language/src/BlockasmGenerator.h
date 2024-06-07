//
// Created by ashy5000 on 6/5/24.
//

#ifndef BLOCKASMGENERATOR_H
#define BLOCKASMGENERATOR_H
#include <string>
#include <vector>

#include "Token.h"
#include "Variable.h"


class BlockasmGenerator {
public:
    std::string GenerateBlockasm(std::vector<Token> tokens);
    std::tuple<std::string, std::vector<Variable>> GenerateSystemFunctionBlockasm(std::vector<Token> tokens, int i, int &nextAllocatedLocation, std::vector<Variable> vars);
};



#endif //BLOCKASMGENERATOR_H
