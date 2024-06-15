//
// Created by ashy5000 on 6/14/24.
//

#ifndef FUNCTION_H
#define FUNCTION_H
#include <string>


class Function {
public:
    int lineOffset;
    std::string name;
    Function(const int lineOffset_p, std::string name_p) : lineOffset(lineOffset_p), name(std::move(name_p)) {}
};



#endif //FUNCTION_H
