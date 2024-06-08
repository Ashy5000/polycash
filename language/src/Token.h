//
// Created by ashy5000 on 6/4/24.
//

#ifndef TOKEN_H
#define TOKEN_H
#include <vector>
#include <string>

#include "TokenType.h"


class Token {
public:
    TokenType type = TokenType::type_placeholder;
    std::string value;
    std::vector<Token> children;
    Token(const TokenType type_p, std::string value_p) {
        type = type_p;
        value = std::move(value_p);
        children = {};
    }
};



#endif //TOKEN_H
