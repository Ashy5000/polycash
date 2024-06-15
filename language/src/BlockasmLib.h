//
// Created by ashy5000 on 6/14/24.
//

#ifndef BLOCKASMLIB_H
#define BLOCKASMLIB_H
#include <string>
#include <vector>

#include "Function.h"


class BlockasmLib {
public:
    std::string sourceFile;
    std::vector<BlockasmLib> dependencies;
    std::vector<Function> functions;
    std::string source;
    void LoadSource();
};



#endif //BLOCKASMLIB_H
