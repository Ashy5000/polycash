//
// Created by ashy5000 on 6/13/24.
//

#include "Signature.h"

#include <utility>


bool Signature::CheckSignature(const std::vector<Type> &types) const {
    for(int i = 0; i < types.size(); i++) {
        if(expectedTypes[i] != types[i] && types[i] != Type::loaded) {
            return false;
        }
    }
    return true;
}

Signature::Signature(std::vector<Type> expectedTypes_p, const Type returnType_p) {
    expectedTypes = std::move(expectedTypes_p);
    returnType = returnType_p;
}

