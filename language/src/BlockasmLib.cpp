//
// Created by ashy5000 on 6/14/24.
//

#include "BlockasmLib.h"

#include <fstream>
#include <sstream>

void BlockasmLib::LoadSource() {
    std::stringstream source_stream;
    {
        std::fstream input("./src/blockasm_lib/" + sourceFile, std::ios::in);
        source_stream << input.rdbuf();
    }
    source = source_stream.str();
    std::istringstream iss(source);
    int i = 0;
    for(std::string line; std::getline(iss, line); ) {
        if(line[0] == ';') {
            std::string name;
            std::istringstream innerIss(line);
            int j = 0;
            for(std::string segment; std::getline(innerIss, segment, '@'); ) {
                if(j == 1) {
                    name = segment;
                    break;
                }
                j++;
            }
            Function f = Function(i + 1, name);
            functions.emplace_back(f);
        }
        i++;
    }
}
