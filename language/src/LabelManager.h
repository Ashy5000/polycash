//
// Created by ashy5000 on 6/24/24.
//

#ifndef POLYCASH_LANG_LABELMANAGER_H
#define POLYCASH_LANG_LABELMANAGER_H

#include <unordered_map>
#include <string>

class LabelManager {
    std::unordered_map<int, int> preLabels;
    std::unordered_map<int, int> labels;
public:
    explicit LabelManager(const std::string &blockasm);
    std::string ReplacePreLabels(const std::string &blockasm);
    std::string ReplaceLabels(const std::string &blockasm);

    static std::string SkipLibs(const std::string &blockasm);
    std::string OffsetCalls(const std::string &blockasm, int offset);
private:
    void DetectPreLabels(const std::string& blockasm);
    void DetectLabels(const std::string& blockasm);
};


#endif //POLYCASH_LANG_LABELMANAGER_H
