//
// Created by ashy5000 on 6/13/24.
//

#ifndef PARAMSPARSER_H
#define PARAMSPARSER_H
#include <tuple>
#include <vector>

#include "Signature.h"
#include "Token.h"
#include "Variable.h"


class ParamsParser {
public:
    std::vector<Token> params;
    Signature sig;
    std::tuple<std::string, std::vector<int>> ParseParams(int &nextAllocatedLocation, const std::vector<Variable>& vars) ;
    ParamsParser(std::vector<Token> params_p, Signature sig_p);
};



#endif //PARAMSPARSER_H
