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
  int selectorLocLocation = -1;

public:
  RegisteredFunctionInfo registerFunction();
  std::string compile(int &nextAllocatedLocation);
  std::string close(int &nextAllocatedLocation);
};



#endif //CONTROLMODULE_HPP
