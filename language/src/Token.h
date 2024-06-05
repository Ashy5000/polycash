//
// Created by ashy5000 on 6/4/24.
//

#ifndef TOKEN_H
#define TOKEN_H
#include <optional>
#include <string>

#include "TokenType.h"


class Token {
    TokenType type = TokenType::type_placeholder;
    std::optional<std::string> value;
public:
    Token(TokenType type_p, std::string value_p) {
        type = type_p;
        value = value_p;
    }
};



#endif //TOKEN_H
