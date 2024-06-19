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
#include "Linker.h"


class ParamsParser {
    std::tuple<std::vector<int>, bool> ParseParamsWithSignature(int &nextAllocatedLocation, const std::vector<Variable> &vars, const Signature& sig, std::stringstream &blockasm, Linker &l);
public:
    std::vector<Token> params;
    std::vector<Signature> signatures;

    std::tuple<std::vector<int>, Signature> ParseParams(int &nextAllocatedLocation,
                                                                     const std::vector<Variable> &vars, std::stringstream &blockasm, Linker &l);
    ParamsParser(std::vector<Token> params_p, std::vector<Signature> signatures_p);
};



#endif //PARAMSPARSER_H
