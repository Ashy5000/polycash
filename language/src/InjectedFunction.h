//
// Created by ashy5000 on 6/16/24.
//

#ifndef INJECTEDFUNCTION_H
#define INJECTEDFUNCTION_H
#include <string>

class InjectedFunction {
public:
    std::string name;
    int offset;

    InjectedFunction(std::string name_p, const int offset_p) : name(std::move(name_p)), offset(offset_p) {}
};

#endif //INJECTEDFUNCTION_H
