//
// Created by ashy5000 on 6/4/24.
//

#ifndef TOKENTYPE_H
#define TOKENTYPE_H



enum class TokenType {
    type_placeholder,
    system_at,
    identifier,
    open_paren,
    close_paren,
    int_lit,
    semi,
    expr,
    newline,
    concat,
    add,
    sub,
    mul,
    div,
    comma,
    eq,
    string_lit,
    excl,
};



#endif //TOKENTYPE_H
