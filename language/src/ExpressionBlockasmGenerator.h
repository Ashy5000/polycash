//
// Created by ashy5000 on 6/5/24.
//

#ifndef EXPRESSIONBLOCKASMGENERATOR_H
#define EXPRESSIONBLOCKASMGENERATOR_H
#include <iostream>
#include <sstream>
#include <string>

#include "Token.h"


class ExpressionBlockasmGenerator {
public:
    std::string GenerateBlockasmFromExpression(Token expression, int nextAllocatedLocation);
};



#endif //EXPRESSIONBLOCKASMGENERATOR_H
