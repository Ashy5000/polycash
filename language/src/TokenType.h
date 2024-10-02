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
    exp,
    comma,
    eq,
    string_lit,
    excl,
    open_curly,
    close_curly,
    greater,
    block,
};



#endif //TOKENTYPE_H
