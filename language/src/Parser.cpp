//
// Created by ashy5000 on 6/4/24.
//

#include "Parser.h"
#include "Token.h"

#include <iostream>

std::vector<Token> Parser::parse_tokens(const std::string &input) {
    std::vector<Token> tokens;
    auto activeToken = Token{TokenType::type_placeholder, {}};
    std::string substring;
    int depth = 0;
    bool inString = false;
    std::string str;
    for(int i = 0; i < input.size(); i++) {
        const char c = input[i];
        if(c == '"' && activeToken.type == TokenType::type_placeholder || activeToken.type == TokenType::string_lit) {
            inString = !inString;
            Token t = Token(TokenType::string_lit, str);
            str.clear();
            if(!inString) {
                tokens.push_back(t);
            }
            continue;
        }
        if(inString) {
            str.push_back(c);
            continue;
        }
        if(c == '{') {
            depth++;
            if(activeToken.type == TokenType::expr || activeToken.type == TokenType::block) {
                substring += "{";
                continue;
            }
            if(!activeToken.value.empty() || (activeToken.type != TokenType::identifier && activeToken.type != TokenType::type_placeholder)) {
                tokens.emplace_back(activeToken);
            }
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            tokens.emplace_back(Token{TokenType::open_curly, {}});
            activeToken.type = TokenType::block;
            continue;
        }
        if(c == '}') {
            if(depth > 1) {
                depth--;
                substring += "}";
                continue;
            }
            depth--;
            activeToken.children = parse_tokens(substring);
            tokens.emplace_back(activeToken);
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            substring.clear();
            tokens.emplace_back(Token{TokenType::close_curly, {}});
        }
        if(c == '(') {
            depth++;
            if(activeToken.type == TokenType::expr || activeToken.type == TokenType::block) {
                substring += "(";
                continue;
            }
            if(!activeToken.value.empty() || (activeToken.type != TokenType::identifier && activeToken.type != TokenType::type_placeholder)) {
                tokens.emplace_back(activeToken);
            }
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            tokens.emplace_back(Token{TokenType::open_paren, {}});
            activeToken.type = TokenType::expr;
            continue;
        }
        if(c == ')') {
            if(depth > 1) {
                depth--;
                substring += ")";
                continue;
            }
            depth--;
            activeToken.children = parse_tokens(substring);
            tokens.emplace_back(activeToken);
            activeToken.children = {};
            activeToken.value = {};
            activeToken.type = TokenType::type_placeholder;
            substring.clear();
            tokens.emplace_back(Token{TokenType::close_paren, {}});
        }
        if(activeToken.type == TokenType::expr || activeToken.type == TokenType::block) {
            substring.push_back(c);
            continue;
        }
        if(c == ' ') {
            continue;
        }
        if(std::isalpha(c) || c == ':' || c == '_' || c == '\'') {
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
        if(c == '!') {
            tokens.emplace_back(Token{TokenType::excl, {}});
        }
        if(c == ';') {
            tokens.emplace_back(Token{TokenType::semi, {}});
        }
        if(c == '.') {
            tokens.emplace_back(Token{TokenType::concat, {}});
        }
        if(c == '+') {
            tokens.emplace_back(Token{TokenType::add, {}});
        }
        if(c == '-') {
            tokens.emplace_back(Token{TokenType::sub, {}});
        }
        if(c == '*') {
            tokens.emplace_back(Token{TokenType::mul, {}});
        }
        if(c == '/') {
            tokens.emplace_back(Token{TokenType::div, {}});
        }
        if(c == '^') {
            tokens.emplace_back(Token{TokenType::exp, {}});
        }
        if(c == ',') {
            tokens.emplace_back(Token{TokenType::comma, {}});
        }
        if(c == '=') {
            tokens.emplace_back(Token{TokenType::eq, {}});
        }
        if(c == '\n') {
            tokens.emplace_back(Token{TokenType::newline, {}});
        }
    }
    return tokens;
}
