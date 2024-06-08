//
// Created by ashy5000 on 6/6/24.
//

#ifndef VARIABLE_H
#define VARIABLE_H
#include <string>
#include <utility>

#include "Type.h"


class Variable {
public:
    std::string name;
    int location;
    Type type;
    Variable(std::string name_p, const int location_p, const Type type_p) {
        name = std::move(name_p);
        location = location_p;
        type = type_p;
    }
};



#endif //VARIABLE_H
