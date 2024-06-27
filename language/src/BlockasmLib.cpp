//
// Created by ashy5000 on 6/14/24.
//

#include "BlockasmLib.h"

#include <fstream>
#include <sstream>
#include <iostream>
#include "Signature.h"

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
            Signature sig({}, Type::type_placeholder);
            int j = 0;
            for(std::string segment; std::getline(innerIss, segment, ' '); ) {
                if(j == 1) {
                    name = segment.substr(1);
                }
                if(j == 2) {
                    std::string returnTypeStr = segment.substr(2);
                    if(returnTypeStr == "string") {
                        sig.returnType = Type::string;
                    } else if(returnTypeStr == "uint64") {
                        sig.returnType = Type::uint64;
                    } else if(returnTypeStr == "boolean") {
                        sig.returnType = Type::boolean;
                    } else {
                        std::cerr << "Unknown type " << returnTypeStr << std::endl;
                        exit(EXIT_FAILURE);
                    }
                }
                if(j > 2) {
                    std::istringstream innerInnerIss(segment);
                    int k = 0;
                    Type paramType;
                    int paramLoc;
                    for(std::string subSegment; std::getline(innerInnerIss, subSegment, '@'); ) {
                        if(k == 0) {
                            std::string paramTypeStr = subSegment.substr(1);
                            if(paramTypeStr == "string") {
                                paramType = Type::string;
                            } else if(paramTypeStr == "uint64") {
                                paramType = Type::uint64;
                            } else {
                                std::cerr << "Unknown type " << paramTypeStr << std::endl;
                                exit(EXIT_FAILURE);
                            }
                        } else if(k == 1) {
                            std::string paramLocStr = subSegment.substr(2);
                            paramLoc = atoi(paramLocStr.c_str());
                        }
                        k++;
                    }
                    sig.expectedTypes.emplace_back(paramType);
                    sig.locations.emplace_back(paramLoc);
                }
                j++;
            }
            Function f = Function(i + 1, name, sig);
            functions.emplace_back(f);
        }
        i++;
    }
}
