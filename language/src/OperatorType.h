//
// Created by ashy5000 on 6/9/24.
//

#ifndef OPERATORTYPE_H
#define OPERATORTYPE_H
#include "TokenType.h"


enum class OperatorType {
    type_placeholder,
    concat,
    add,
    sub,
    mul,
    div,
    eq,
    exp
};

inline OperatorType OperatorTypeFromToken(const Token& t) {
    switch(t.type) {
        case TokenType::concat:
            return OperatorType::concat;
        case TokenType::add:
            return OperatorType::add;
        case TokenType::sub:
            return OperatorType::sub;
        case TokenType::mul:
            return OperatorType::mul;
        case TokenType::div:
            return OperatorType::div;
        case TokenType::eq:
            return OperatorType::eq;
        case TokenType::exp:
            return OperatorType::exp;
        default:
            return OperatorType::type_placeholder;
    }
}

#endif //OPERATORTYPE_H
