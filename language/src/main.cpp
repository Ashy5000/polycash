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
        std::cerr << "plc [in.pcon]" << std::endl;
    }

    std::stringstream contents_stream;
    {
        std::fstream input(argv[1], std::ios::in);
        contents_stream << input.rdbuf();
    }
    std::string contents = contents_stream.str();
    Parser parser;
    std::vector<Token> tokens = parser.parse_tokens(contents);
    BlockasmGenerator generator = BlockasmGenerator(tokens);
    std::string blockasm = generator.GenerateBlockasm();
    std::ofstream targetAsm("out.blockasm");
    if (targetAsm.is_open()) {
        targetAsm << blockasm << std::endl;
        targetAsm.close();
    } else {
        std::cerr << "Error opening target file." << std::endl;
    }
    return 0;
}
