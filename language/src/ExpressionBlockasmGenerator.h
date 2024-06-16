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
    static std::tuple<std::string, int, Type> GenerateBlockasmFromExpression(const Token &expression, int nextAllocatedLocation, const std::vector<Variable> &vars);
};



#endif //EXPRESSIONBLOCKASMGENERATOR_H
