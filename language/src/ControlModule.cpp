//
// Created by ashy5000 on 10/6/24.
//

#include "ControlModule.hpp"

#include <algorithm>
#include <iomanip>
#include <iostream>
#include <sstream>

RegisteredFunctionInfo ControlModule::registerFunction() {
  RegisteredFunctionInfo info;
  info.preLabelId = static_cast<int>(registeredFunctionInfos.size());
  info.id = rand();
  registeredFunctionInfos.emplace_back(info);
  return info;
}

std::string ControlModule::compile(int &nextAllocatedLocation) {
  std::stringstream result;
  // Set selector location
  int selectorLocLocation = nextAllocatedLocation++;
  result << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << selectorLocLocation << " 0x00000000" << std::endl;
  result << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << selectorLocLocation << " 0x010000 0x00000000" << std::endl;
  // Get selected function ID
  int selectorLocation = nextAllocatedLocation++;
  result << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << selectorLocation << " 0x00000000" << std::endl;
  result << "ReadFromState 0x" << std::setfill('0') << std::setw(8) << std::hex << selectorLocLocation << " 0x";
  result << std::setfill('0') << std::setw(8) << std::hex << selectorLocation << " 0x00000000" << std::endl;
  // Create buffer to hold comparison result
  int cmpResLocation = nextAllocatedLocation++;
  result << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << cmpResLocation << " 0x00000000" << std::endl;
  // Create buffer to hold current ID
  int idLocation = nextAllocatedLocation++;
  result << "InitBfr 0x" << std::setfill('0') << std::setw(8) << std::hex << idLocation << " 0x00000000" << std::endl;
  for(RegisteredFunctionInfo info : registeredFunctionInfos) {
    // Set current ID
    result << "SetCnst 0x" << std::setfill('0') << std::setw(8) << std::hex << idLocation << " 0x";
    result << std::setfill('0') << std::setw(8) << std::hex << info.id << " 0x00000000" << std::endl;
    // Compare
    result << "Eq 0x" << std::setfill('0') << std::setw(8) << std::hex << idLocation << " 0x";
    result << std::setfill('0') << std::setw(8) << std::hex << selectorLocation << " 0x";
    result << std::setfill('0') << std::setw(8) << std::hex << cmpResLocation << " 0x00000000" << std::endl;
    // Jump
    result << "JmpCond 0x" << std::setfill('0') << std::setw(8) << std::hex << cmpResLocation << " !";
    result << info.preLabelId << " 0x00000000" << std::endl;
  }
  return result.str();
}

