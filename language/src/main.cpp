//
// Created by ashy5000 on 6/4/24.
//

#include <fstream>
#include <iostream>
#include <sstream>
#include <string>

#include "BlockasmGenerator.h"
#include "Parser.h"

int main(int argc, char* argv[]) {
    if (argc != 2) {
        std::cerr << "Incorrect usage." << std::endl;
        std::cerr << "plc [in.poly]" << std::endl;
    }

    std::stringstream contents_stream;
    {
        std::fstream input(argv[1], std::ios::in);
        contents_stream << input.rdbuf();
    }
    std::string contents = contents_stream.str();
    std::vector<Token> tokens = Parser::parse_tokens(contents);
    auto generator = BlockasmGenerator(tokens, 0x00001000, {}, true);
    std::string blockasm = generator.GenerateBlockasm();
    if (std::ofstream targetAsm("out.blockasm"); targetAsm.is_open()) {
        targetAsm << blockasm << std::endl;
        targetAsm.close();
    } else {
        std::cerr << "Error opening target file." << std::endl;
    }
    return 0;
}
