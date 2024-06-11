//
// Created by ashy5000 on 6/5/24.
//

#ifndef BLOCKASMGENERATOR_H
#define BLOCKASMGENERATOR_H
#include <sstream>
#include <string>
#include <vector>

#include "Token.h"
#include "Variable.h"


class BlockasmGenerator {
public:
    std::string GenerateBlockasm();
    std::tuple<std::vector<Variable>, int> GenerateSystemFunctionBlockasm(int i, int &nextAllocatedLocation, std::vector<Variable> vars);
    explicit BlockasmGenerator(std::vector<Token> tokens_p);
private:
    std::stringstream blockasm;
    std::vector<Token> tokens;
};



#endif //BLOCKASMGENERATOR_H
