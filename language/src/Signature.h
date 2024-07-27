//
// Created by ashy5000 on 6/13/24.
//

#ifndef SIGNATURE_H
#define SIGNATURE_H
#include <vector>

#include "Type.h"


class Signature {
public:
    std::vector<Type> expectedTypes;
    std::vector<int> locations;
    Type returnType;
    [[nodiscard]] bool CheckSignature(const std::vector<Type> &types) const;
    explicit Signature(std::vector<Type> expectedTypes_p, Type returnType_p);
};



#endif //SIGNATURE_H
