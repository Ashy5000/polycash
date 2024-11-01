//
// Created by ashy5000 on 6/9/24.
//

#ifndef OPERATOR_H
#define OPERATOR_H
#include "OperatorType.h"
#include "Type.h"
#include <string>

class Operator {
public:
    OperatorType type;
};

inline std::string OperatorToString(const Operator o) {
    switch(o.type) {
        case OperatorType::type_placeholder:
            return "";
        case OperatorType::concat:
            return "App";
        case OperatorType::add:
            return "Add";
        case OperatorType::sub:
            return "Sub";
        case OperatorType::mul:
            return "Mul";
        case OperatorType::div:
            return "Div";
        case OperatorType::exp:
            return "Exp";
        case OperatorType::eq:
            return "Eq";
    }
    return "";
}

inline Type OperatorToType(const Operator o) {
    switch(o.type) {
        case OperatorType::type_placeholder:
            return Type::type_placeholder;
        case OperatorType::add:
        case OperatorType::sub:
        case OperatorType::mul:
        case OperatorType::div:
        case OperatorType::exp:
            return Type::uint64;
        case OperatorType::concat:
            return Type::string;
        case OperatorType::eq:
            return Type::boolean;
    }
    return Type::type_placeholder;
}


#endif //OPERATOR_H
