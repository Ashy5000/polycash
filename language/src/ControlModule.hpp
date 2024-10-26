//
// Created by ashy5000 on 10/6/24.
//

#ifndef CONTROLMODULE_HPP
#define CONTROLMODULE_HPP

#include <string>
#include <vector>

#include "RegisteredFunctionInfo.h"


class ControlModule {
  std::vector<RegisteredFunctionInfo> registeredFunctionInfos;

public:
  RegisteredFunctionInfo registerFunction();
  std::string compile(int &nextAllocatedLocation);
};



#endif //CONTROLMODULE_HPP
