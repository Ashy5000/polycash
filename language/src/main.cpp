//
// Created by ashy5000 on 6/4/24.
//

#include <fstream>
#include <iostream>
#include <string>

#include "Compiler.h"

int main(int argc, char* argv[]) {
    if (argc != 2) {
    std::cerr << "Incorrect usage." << std::endl;
        std::cerr << "plc [in.poly]" << std::endl;
    }

    srand(time(nullptr));

    auto c = Compiler();
    std::string blockasm = c.Compile(argv[1]);
    if (std::ofstream targetAsm("out.blockasm"); targetAsm.is_open()) {
        targetAsm << blockasm << std::endl;
        targetAsm.close();
    } else {
        std::cerr << "Error opening target file." << std::endl;
    }
    return 0;
}
