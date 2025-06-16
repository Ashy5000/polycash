//
// Created by ashy5000 on 6/13/24.
//

#ifndef SYSTEMFUNCTION_H
#define SYSTEMFUNCTION_H
#include <functional>
#include <string>

#include "Token.h"
#include "Variable.h"
#include "Linker.h"


class SystemFunction {
public:
    std::function<void(std::vector<Token>, int&, std::vector<Variable>&, std::stringstream&, Linker&)> generateBlockasm;
    std::string module;
    std::string name;
    SystemFunction(std::function<void(std::vector<Token>, int&, std::vector<Variable>&, std::stringstream&, Linker&)> generateBlockasm_p, std::string module_p, std::string name_p) : generateBlockasm(generateBlockasm_p), module(std::move(module_p)), name(std::move(name_p)) {}
};



#endif //SYSTEMFUNCTION_H
