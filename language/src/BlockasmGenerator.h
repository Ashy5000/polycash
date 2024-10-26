//
// Created by ashy5000 on 6/5/24.
//

#ifndef BLOCKASMGENERATOR_H
#define BLOCKASMGENERATOR_H
#include <sstream>
#include <string>
#include <tuple>
#include <vector>

#include "ControlModule.hpp"
#include "Linker.h"
#include "Token.h"
#include "Variable.h"


class BlockasmGenerator {
public:
    std::string GenerateBlockasm(ControlModule cm);
    std::tuple<std::vector<Variable>, int> GenerateSystemFunctionBlockasm(int i, int &nextAllocatedLocation, std::vector<Variable> vars, Linker l);
    explicit BlockasmGenerator(std::vector<Token> tokens_p, int nextAllocatedLocation_p, std::vector<Variable> vars_p, bool useLinker_p);
    int GetNextAllocatedLocation() const;
private:
    std::stringstream blockasm;
    std::vector<Token> tokens;
    std::vector<Variable> vars;
    int nextAllocatedLocation;
    int nextAllocatedStateLocation;
    bool useLinker;
};



#endif //BLOCKASMGENERATOR_H
