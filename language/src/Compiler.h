#ifndef COMPILER_H
#define COMPILER_H
#include <string>
#include <vector>
#include "Token.h"

class Compiler {
    std::string filename;
    std::string contents;
    std::vector<Token> tokens;
    std::string blockasm;
public:
    void LoadContents();
    void ParseTokens();
    void GenerateBlockasm();
    void Link();
    std::string Compile(std::string filename);
};

#endif
