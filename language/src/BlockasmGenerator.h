//
// Created by ashy5000 on 6/5/24.
//

#ifndef BLOCKASMGENERATOR_H
#define BLOCKASMGENERATOR_H
#include <string>
#include <vector>

#include "Token.h"


class BlockasmGenerator {
public:
    std::string GenerateBlockasm(std::vector<Token> tokens);
    std::string GenerateSystemFunctionBlockasm(std::vector<Token> tokens, int i, int &nextAllocatedLocation);
};



#endif //BLOCKASMGENERATOR_H
