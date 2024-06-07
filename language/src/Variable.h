//
// Created by ashy5000 on 6/6/24.
//

#ifndef VARIABLE_H
#define VARIABLE_H
#include <string>


class Variable {
public:
    std::string name;
    int location;
    Variable(std::string name_p, int location_p) {
        name = name_p;
        location = location_p;
    }
};



#endif //VARIABLE_H
