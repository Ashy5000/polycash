//
// Created by ashy5000 on 6/4/24.
//

#include <fstream>
#include <iostream>
#include <sstream>
#include <string>

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
    return 0;
}
