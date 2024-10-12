//
// Created by ashy5000 on 10/6/24.
//

#include "ControlModule.hpp"

#include <algorithm>

// Adapted from https://stackoverflow.com/a/45688206
std::string generateUniqueString() {
  auto randchar = []() -> char
  {
    const char charset[] = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
    const size_t max_index = (sizeof(charset) - 1);
    return charset[ rand() % max_index ];
  };
  std::string str(32,0);
  std::generate_n( str.begin(), 4, randchar );
  return str;
}

RegisteredFunctionInfo ControlModule::registerFunction() {
  RegisteredFunctionInfo info;
  info.preLabelId = static_cast<int>(registeredFunctionInfos.size());
  info.id = generateUniqueString();
  registeredFunctionInfos.emplace_back(info);
  return info;
}
