//
// Created by ashy5000 on 6/24/24.
//

#ifndef POLYCASH_LANG_LABELMANAGER_H
#define POLYCASH_LANG_LABELMANAGER_H

#include <unordered_map>
#include <string>

class LabelManager {
    std::unordered_map<int, int> labels;
public:
    explicit LabelManager(std::string blockasm);
    std::string ReplaceLabels(std::string blockasm);
};


#endif //POLYCASH_LANG_LABELMANAGER_H
