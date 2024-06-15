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

const std::vector SYSTEM_FUNCTIONS = {
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, const std::vector<Variable>& vars) -> std::string {
            std::stringstream blockasm;
            Signature sig = Signature({Type::uint64});
            ParamsParser pp = ParamsParser(params, {sig});
            std::tuple<std::string, std::vector<int>, Signature> parsingResult = pp.ParseParams(nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(parsingResult);
            std::vector<int> locations = std::get<1>(parsingResult);
            int exitCodeLocation = locations[0];
            blockasm << expressionBlockasm;
            blockasm << "ExitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << exitCodeLocation;
            if(exitCodeLocation >= nextAllocatedLocation) {
                nextAllocatedLocation = exitCodeLocation + 1;
            }
            return blockasm.str();
        },
        "contract",
        "exit"
    ),
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, std::vector<Variable> &vars) -> std::string {
            std::stringstream blockasm;
            std::string typeString = params[1].children[0].value;
            Type type = Type::type_placeholder;
            if(typeString == "uint64") {
                type = Type::uint64;
            }
            if(type == Type::type_placeholder) {
                std::cerr << "Unknown type " << typeString << std::endl;
                exit(EXIT_FAILURE);
            }
            auto var = Variable(params[0].children[0].value, nextAllocatedLocation, type);
            vars.emplace_back(var);
            blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << nextAllocatedLocation++ << " 0x00000000" << std::endl;
            return blockasm.str();
        },
        "memory",
        "alloc"
    ),
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, std::vector<Variable> &vars) -> std::string {
            std::stringstream blockasm;
            int indexToRemove = -1;
            for(int j = 0; j < vars.size(); j++) {
                if(Variable var = vars[j]; var.name == params[0].children[0].value) {
                    indexToRemove = j;
                    break;
                }
            }
            if(indexToRemove == -1) {
                std::cerr << "Cannot free undefined variable." << std::endl;
                exit(EXIT_FAILURE);
            }
            blockasm << "Free 0x" << std::setfill('0') << std::setw(8) << std::hex << vars[indexToRemove].location << " 0x00000000" << std::endl;
            vars.erase(vars.begin() + indexToRemove);
            return blockasm.str();
        },
        "memory",
        "free"
    ),
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, std::vector<Variable> &vars) -> std::string {
            std::stringstream blockasm;
            int indexToRename = -1;
            for(int j = 0; j < vars.size(); j++) {
                if(Variable var = vars[j]; var.name == params[0].children[0].value) {
                    indexToRename = j;
                    break;
                }
            }
            if(indexToRename == -1) {
                std::cerr << "Cannot set undefined variable." << std::endl;
                exit(EXIT_FAILURE);
            }
            char* end;
            int val = static_cast<int>(std::strtol(params[1].children[0].value.c_str(), &end, 10));
            if(errno == ERANGE) {
                std::cerr << "Expected integer as value" << std::endl;
            }
            blockasm << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << vars[indexToRename].location << " 0x";
            blockasm << std::setfill('0') << std::setw(16) << std::hex << val << " 0x00000000" << std::endl;
            return blockasm.str();
        },
        "memory",
        "set"
    ),
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, std::vector<Variable> &vars) -> std::string {
            std::stringstream blockasm;
            Signature uint64Sig = Signature({Type::uint64});
            Signature stringSig = Signature({Type::string});
            ParamsParser pp = ParamsParser(params, {uint64Sig, stringSig});
            std::tuple<std::string, std::vector<int>, Signature> parsingResult = pp.ParseParams(nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(parsingResult);
            std::vector<int> locations = std::get<1>(parsingResult);
            Signature sig = std::get<2>(parsingResult);
            blockasm << expressionBlockasm << std::endl;
            int dataLocation = locations[0];
            if(sig.expectedTypes[0] == Type::uint64) {
                blockasm << "Stdout 0x" << std::setfill('0') << std::setw(8) << std::hex << dataLocation << " 0x00000000" << std::endl;
            } else {
                blockasm << "PrintStr 0x" << std::setfill('0') << std::setw(8) << std::hex << dataLocation << " 0x00000000" << std::endl;
            }
            return blockasm.str();
        },
        "io",
        "print"
    ),
    SystemFunction(
        [](const std::vector<Token>& params, int &nextAllocatedLocation, std::vector<Variable> &vars) -> std::string {
            std::stringstream blockasm;
            Signature sig = Signature({Type::uint64});
            ParamsParser pp = ParamsParser(params, {sig});
            std::tuple<std::string, std::vector<int>, Signature> parsingResult = pp.ParseParams(nextAllocatedLocation, vars);
            std::string expressionBlockasm = std::get<0>(parsingResult);
            std::vector<int> locations = std::get<1>(parsingResult);
            blockasm << expressionBlockasm << std::endl;
            int dataLocation = locations[0];
            blockasm << "Stderr 0x" << std::setfill('0') << std::setw(8) << std::hex << dataLocation << " 0x00000000" << std::endl;
            return blockasm.str();
        },
        "io",
        "err"
    )
};

#endif //SYSTEMFUNCTIONS_H
