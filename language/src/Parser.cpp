//
// Created by ashy5000 on 6/4/24.
//

#include "Parser.h"
#include "Token.h"

#include <iostream>

std::vector<Token> Parser::parse_tokens(std::string input) {
    std::vector<Token> tokens;
    std::vector<char> buf;
    TokenType currentType;
    for(char c : input) {
        if(std::isalpha(c)) {
            if(buf.empty()) {
                currentType = TokenType::identifier;
            }
            if(currentType != TokenType::identifier) {
                std::cerr << "Unexpected token " << c << "." << std::endl;
                exit(EXIT_FAILURE);
            }
            buf.push_back(c);
            continue;
        }
        if(std::isdigit(c)) {
            if(buf.empty()) {
                currentType = TokenType::int_lit;
            }
            buf.push_back(c);
            continue;
        }
        if(!buf.empty()) {
            std::string str(buf.begin(), buf.end());
            tokens.emplace_back(currentType, str);
            buf.clear();
        }
        if(std::isspace(c)) {
            continue;
        }
        if(c == '@') {
            tokens.emplace_back(Token{TokenType::system_at, {}});
            continue;
        }
        if(c == ';') {
            tokens.emplace_back(Token{TokenType::semi, {}});
            continue;
        }
        if(c == '(') {
            tokens.emplace_back(Token{TokenType::open_paren, {}});
            continue;
        }
        if(c == ')') {
            tokens.emplace_back(Token{TokenType::close_paren, {}});
        }
    }
    return tokens;
}
