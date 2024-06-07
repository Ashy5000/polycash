//
// Created by ashy5000 on 6/5/24.
//

#ifndef EXPRESSIONBLOCKASMGENERATOR_H
#define EXPRESSIONBLOCKASMGENERATOR_H
#include <string>

#include "Token.h"
#include "Variable.h"


class ExpressionBlockasmGenerator {
public:
    std::tuple<std::string, int> GenerateBlockasmFromExpression(Token expression, int nextAllocatedLocation, std::vector<Variable> vars);
};



#endif //EXPRESSIONBLOCKASMGENERATOR_H
