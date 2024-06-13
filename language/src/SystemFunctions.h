//
// Created by ashy5000 on 6/13/24.
//

#ifndef SYSTEMFUNCTIONS_H
#define SYSTEMFUNCTIONS_H
#include <iomanip>
#include <sstream>

#include "ParamsParser.h"
#include "Signature.h"
#include "SystemFunction.h"

const std::vector<SystemFunction> SYSTEM_FUNCTIONS = {
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, const std::vector<Variable>& vars) -> std::string {
            std::stringstream blockasm;
            Signature sig = Signature({Type::uint64});
            ParamsParser pp = ParamsParser(params, sig);
            std::tuple<std::string, std::vector<int>> parsingResult = pp.ParseParams(nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(parsingResult);
            std::vector<int> locations = std::get<1>(parsingResult);
            int exitCodeLocation = locations[0];
            blockasm << expressionBlockasm;
            blockasm << "ExitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << exitCodeLocation << std::endl;
            if(exitCodeLocation >= nextAllocatedLocation) {
                nextAllocatedLocation = exitCodeLocation + 1;
            }
            return blockasm.str();
        },
        "contract",
        "exit"
    )
};

#endif //SYSTEMFUNCTIONS_H
