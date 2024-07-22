//
// Created by ashy5000 on 6/5/24.
//

#ifndef EXPRESSIONBLOCKASMGENERATOR_H
#define EXPRESSIONBLOCKASMGENERATOR_H
#include <string>
#include <tuple>

#include "Token.h"
#include "Variable.h"
#include "Linker.h"


class ExpressionBlockasmGenerator {
public:
    static std::tuple<int, Type> GenerateBlockasmFromExpression(Token expression, int nextAllocatedLocation, std::vector<Variable> &vars, std::stringstream &blockasm, Linker &l);
};



#endif //EXPRESSIONBLOCKASMGENERATOR_H
