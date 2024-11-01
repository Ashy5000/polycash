//
// Created by ashy5000 on 6/14/24.
//

#include "Linker.h"

#include <iostream>
#include <iomanip>
#include <sstream>
#include <tuple>

#include "Variable.h"

Linker::Linker(const std::vector<std::string> &entries) {
    libs = {};
    functionsInjected = {};
    for(const auto &entry : entries) {
        auto lib = BlockasmLib();
        lib.sourceFile = entry;
        lib.LoadSource();
        libs.emplace_back(lib);
    }
}

void Linker::InjectIfNotPresent(const std::string& name, std::stringstream &blockasm) {
    for(const InjectedFunction& fn: functionsInjected) {
        if(fn.name == name) {
            return;
        }
    }
    int offset = -1;
    for(const auto&[sourceFile, dependencies, functions, source] : libs) {
        for(const auto& func : functions) {
            if(func.name == name) {
                offset = 2;
                std::string result;
                std::istringstream iss(result);
                int j = 0;
                std::stringstream before;
                std::stringstream after;
                bool streamingToBefore = true;
                for(std::string line; std::getline(blockasm, line); ) {
                    if(line == ";^^^^BEGIN_SOURCE^^^^") {
                        after << line << std::endl;
                        streamingToBefore = false;
                        continue;
                    }
                    if(streamingToBefore) {
                        offset++;
                        before << line << std::endl;
                        continue;
                    }
                    after << line << std::endl;
                }
                for(std::string line; std::getline(iss, line); ) {
                    if(j >= func.lineOffset) {
                        if(line.substr(0, 3) == "Jmp") {
                            std::stringstream ss(line);
                            std::string temp;
                            std::stringstream adjustedLine;
                            while(ss >> temp) {
                                if(temp.substr(0, 2) != "0x" && temp.substr(0, 3) != "Jmp") {
                                    std::string tempClean;
                                    bool cleaned = false;
                                    if(temp.at(0) == '&') {
                                        tempClean = temp.substr(1);
                                        cleaned = true;
                                    } else {
                                        tempClean = temp;
                                    }
                                    int relativeLine = stoi(tempClean);
                                    int absoluteLine = relativeLine + offset - 1;
                                    if(cleaned) {
                                        adjustedLine << "&" << absoluteLine;
                                    } else {
                                        adjustedLine << absoluteLine;
                                    }
                                } else {
                                    adjustedLine << temp;
                                }
                                adjustedLine << " ";
                            }
                            result += adjustedLine.str() + "\n";
                        } else if(line.at(0) == '%') {
                            std::string functionName = line.substr(1);
                            InjectIfNotPresent(functionName, blockasm);
                            int functionOffset = -1;
                            for(const InjectedFunction& injected_function : functionsInjected) {
                                if(injected_function.name == functionName) {
                                    functionOffset = injected_function.offset;
                                    break;
                                }
                            }
                            if(functionOffset == -1) {
                                std::cerr << "Function not found!" << std::endl;
                                exit(EXIT_FAILURE);
                            }
                            result += "Call &";
                            result += std::to_string(functionOffset);
                            result += "\n";
                        } else {
                            result += line + "\n";
                        }
                        if(line.substr(0, 3) == "Ret") {
                            break;
                        }
                    }
                    j++;
                }
                blockasm = std::stringstream();
                blockasm << before.str();
                blockasm << result;
                std::string blockasmStr = blockasm.str();
                blockasm << after.str();
            }
        }
    }
    if(offset == -1) {
        std::cerr << "Unknown function " << name << std::endl;
        exit(EXIT_FAILURE);
    }
    auto injectedFunction = InjectedFunction(name, offset);
    functionsInjected.emplace_back(injectedFunction);
}

void Linker::SkipLibs(std::stringstream &blockasm) {
    const std::string &temp = blockasm.str();
    blockasm = {};
    blockasm << "Jmp %" << std::endl;
    blockasm << temp;
}

std::tuple<std::string, Type> Linker::CallFunction(const std::string& name, const std::vector<int> &paramLocs, std::vector<Variable>& vars) {
    Type t;
    for(const InjectedFunction& func : functionsInjected) {
        if(func.name == name) {
            std::stringstream blockasm;
            for(int i = 0; i < paramLocs.size(); i++) {
                const int fromLoc = paramLocs[i];
                int toLoc = 0;
                for(const auto&[sourceFile, dependencies, functions, source] : libs) {
                    for(Function libFunc : functions) {
                        if(libFunc.name == func.name) {
                            toLoc = libFunc.sig.locations[i];
                            t = libFunc.sig.returnType;
                            break;
                        }
                    }
                }
                bool exists = false;
                for(const Variable& var : vars) {
                    if(var.location == toLoc) {
                        exists = true;
                    }
                }
                if(!exists) {
                    vars.emplace_back("", toLoc, Type::type_placeholder);
                    blockasm << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << toLoc << " 0x00000000" << std::endl;
                }
                blockasm << "CpyBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << fromLoc << " 0x";
                blockasm << std::setfill('0') << std::setw(8) << std::hex << toLoc << " 0x00000000" << std::endl;
            }
            blockasm << "Call &" << std::dec << func.offset << std::endl;
            return std::make_tuple(blockasm.str(), t);
        }
    }
    std::cerr << "Attempt to call unresolved function " << name << std::endl;
    exit(EXIT_FAILURE);
}

