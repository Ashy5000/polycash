//
// Created by ashy5000 on 6/14/24.
//

#include "Linker.h"

#include <iostream>
#include <iomanip>
#include <sstream>

Linker::Linker(const std::vector<std::string> &entries) {
    libs = {};
    functionsInjected = {};
    for(const auto &entry : entries) {
        BlockasmLib lib = BlockasmLib();
        lib.sourceFile = entry;
        lib.LoadSource();
        libs.emplace_back(lib);
    }
}

void Linker::InjectIfNotPresent(std::string name, std::stringstream &blockasm) {
    for(const InjectedFunction& fn: functionsInjected) {
        if(fn.name == name) {
            return;
        }
    }
    int offset = -1;
    for(const BlockasmLib& lib : libs) {
        for(int i = 0; i < lib.functions.size(); i++) {
            Function func = lib.functions[i];
            if(func.name == name) {
                offset = 2;
                std::string source;
                std::istringstream iss(lib.source);
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
                                    int relativeLine = stoi(temp);
                                    int absoluteLine = relativeLine + offset - 1;
                                    adjustedLine << absoluteLine;
                                } else {
                                    adjustedLine << temp;
                                }
                                adjustedLine << " ";
                            }
                            source += adjustedLine.str() + "\n";
                        } else {
                            source += line + "\n";
                        }
                    }
                    j++;
                }
                blockasm = std::stringstream();
                blockasm << before.str();
                blockasm << source;
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
    int jmpTo = -1;
    int i = 0;
    std::istringstream iss(blockasm.str());
    for(std::string line; std::getline(iss, line); ) {
        if(line == ";^^^^BEGIN_SOURCE^^^^") {
            jmpTo = i + 2;
            break;
        }
        i++;
    }
    if(jmpTo == -1) {
        std::cerr << "Could not skip libraries; no ^^^^BEGIN_SOURCE^^^^ declaration found." << std::endl;
        exit(EXIT_FAILURE);
    }
    const std::string &temp = blockasm.str();
    blockasm.seekp(0);
    blockasm << "Jmp " << jmpTo << std::endl;
    blockasm << temp;
}

std::string Linker::CallFunction(const std::string& name, std::vector<int> paramLocs) {
    for(const InjectedFunction& func : functionsInjected) {
        if(func.name == name) {
            std::stringstream blockasm;
            for(int i = 0; i < paramLocs.size(); i++) {
                int fromLoc = paramLocs[i];
                int toLoc = 0;
                for(BlockasmLib lib : libs) {
                    for(Function libFunc : lib.functions) {
                        if(libFunc.name == func.name) {
                            toLoc = libFunc.sig.locations[i];
                            break;
                        }
                    }
                }
                blockasm << "CpyBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << fromLoc << " 0x";
                blockasm << std::setfill('0') << std::setw(8) << std::hex << toLoc << " 0x00000000" << std::endl;
            }
            blockasm << "Call " << func.offset << std::endl;
            return blockasm.str();
        }
    }
    std::cerr << "Attempt to call unresolved function " << name << std::endl;
    exit(EXIT_FAILURE);
}

