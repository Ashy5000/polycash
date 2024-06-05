//
// Created by ashy5000 on 6/4/24.
//

#include "Parser.h"
#include "Token.h"

#include <iostream>

std::vector<Token> Parser::parse_tokens(std::string input) {
    std::vector<Token> tokens;
    Token activeToken = Token{TokenType::type_placeholder, {}};
    std::string substring;
    for(int i = 0; i < input.size(); i++) {
        char c = input[i];
        if(c == '(') {
            tokens.emplace_back(activeToken);
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            tokens.emplace_back(Token{TokenType::open_paren, {}});
            activeToken.type = TokenType::expr;
            continue;
        }
        if(c == ')') {
            activeToken.children = parse_tokens(substring);
            tokens.emplace_back(activeToken);
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            substring.clear();
            tokens.emplace_back(Token{TokenType::close_paren, {}});
        }
        if(activeToken.type == TokenType::expr) {
            substring.push_back(c);
            continue;
        }
        if(c == ' ') {
            continue;
        }
        if(std::isalpha(c)) {
            if(activeToken.value.empty()) {
                activeToken.type = TokenType::identifier;
            }
            if(activeToken.type != TokenType::identifier) {
                std::cerr << "Unexpected token " << c << "." << std::endl;
                exit(EXIT_FAILURE);
            }
            activeToken.value.push_back(c);
            if(i != input.size() - 1) {
                continue;
            }
        }
        if(std::isdigit(c)) {
            if(activeToken.value.empty()) {
                activeToken.type = TokenType::int_lit;
            }
            activeToken.value.push_back(c);
            if(i != input.size() - 1) {
                continue;
            }
        }
        if(!activeToken.value.empty()) {
            std::string str(activeToken.value.begin(), activeToken.value.end());
            tokens.emplace_back(activeToken.type, str);
            activeToken.value.clear();
        }
        if(c == '@') {
            tokens.emplace_back(Token{TokenType::system_at, {}});
        }
        if(c == ';') {
            tokens.emplace_back(Token{TokenType::semi, {}});
        }
    }
    return tokens;
}
