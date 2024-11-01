#ifndef COMPILER_H
#define COMPILER_H
#include <string>
#include <vector>
#include <bits/random.h>

#include "Token.h"

class Compiler {
    std::string filename;
    std::string contents;
    std::vector<Token> tokens;
    std::string blockasm;
    int controlSegmentSize = 0;
public:
    void LoadContents();
    void ParseTokens();
    void GenerateBlockasm(std::default_random_engine rnd);
    void Link();
    std::string Compile(const std::string &filename_p, std::default_random_engine rnd);
};

#endif
