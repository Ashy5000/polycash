//
// Created by ashy5000 on 6/5/24.
//

#ifndef EXPRESSIONBLOCKASMGENERATOR_H
#define EXPRESSIONBLOCKASMGENERATOR_H
#include <string>
#include <tuple>

#include "Token.h"
#include "Variable.h"


class ExpressionBlockasmGenerator {
public:
    static std::tuple<std::string, int> GenerateBlockasmFromExpression(const Token &expression, const int nextAllocatedLocation, const std::vector<Variable> &vars);
};



#endif //EXPRESSIONBLOCKASMGENERATOR_H
