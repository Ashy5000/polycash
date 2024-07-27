//
// Created by ashy5000 on 6/14/24.
//

#ifndef LINKER_H
#define LINKER_H
#include "BlockasmLib.h"
#include "InjectedFunction.h"
#include "Variable.h"


class Linker {
    std::vector<InjectedFunction> functionsInjected;
public:
    std::vector<BlockasmLib> libs;
    void InjectIfNotPresent(const std::string& name, std::stringstream &blockasm);
    [[nodiscard]] std::tuple<std::string, Type> CallFunction(const std::string &name, std::vector<int> paramLocs, std::vector<Variable>& vars);

    static void SkipLibs(std::stringstream &blockasm);

    explicit Linker(const std::vector<std::string> &entries);
};



#endif //LINKER_H
