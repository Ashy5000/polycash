#ifndef COMPILER_H
#define COMPILER_H
#include <string>
#include <vector>
#include <random>

#include "Token.h"
#include "ControlModule.hpp"

class Compiler {
    std::string filename;
    std::string contents;
    std::vector<Token> tokens;
    std::string blockasm;
    int controlSegmentSize = 0;
    void LoadContents();
    void ParseTokens();
    void InjectControlBuffer(ControlModule &cm);
    void GenerateBlockasm(std::default_random_engine rnd);
    void Link();
    void NullifyString(const std::string& string);
    void Cleanup();
public:
    std::string Compile(const std::string &filename_p, std::default_random_engine rnd);
};

#endif
