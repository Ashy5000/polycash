//
// Created by ashy5000 on 6/14/24.
//

#ifndef LINKER_H
#define LINKER_H
#include "BlockasmLib.h"


class Linker {
    std::vector<std::string> functionsInjected;
public:
    std::vector<BlockasmLib> libs;
    void InjectIfNotPresent(std::string name, std::stringstream &blockasm);

    static void SkipLibs(std::stringstream &blockasm);

    explicit Linker(const std::vector<std::string> &entries);
};



#endif //LINKER_H
