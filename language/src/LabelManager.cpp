//
// Created by ashy5000 on 6/24/24.
//

#include <sstream>
#include "LabelManager.h"

LabelManager::LabelManager(std::string blockasm) {
    labels = {};
    std::istringstream iss(blockasm);
    int labelPosition = 0;
    for (std::string line; std::getline(iss, line); ) {
        if(line.size() < 1 || line.at(0) != ';') {
            labelPosition++;
            continue;
        }
        if(line.substr(2, 5) != "LABEL") {
            labelPosition++;
            continue;
        }
        std::string labelIdString = line.substr(8);
        int labelId = stoi(labelIdString);
        labels[labelId] = labelPosition;
        labelPosition++;
    }
}

std::string LabelManager::ReplaceLabels(std::string blockasm) {
    std::istringstream iss(blockasm);
    std::stringstream newBlockasm;
    for(std::string line; std::getline(iss, line); ) {
        if(line.substr(0, 3) != "Jmp") {
            newBlockasm << line << std::endl;
            continue;
        }
        std::istringstream innerIss(line);
        std::string section;
        std::stringstream newLine;
        bool isFirstSection = true;
        while(innerIss >> section) {
            if(section.at(0) != '<') {
                if(isFirstSection) {
                    newLine << section;
                    isFirstSection = false;
                    continue;
                }
                newLine << " " << section;
                continue;
            }
            std::string labelIdString = section.substr(1);
            int labelId = stoi(labelIdString);
            int labelPos = labels[labelId];
            if(isFirstSection) {
                newLine << labelPos;
                isFirstSection = false;
                continue;
            }
            newLine << " " << labelPos;
        }
        newBlockasm << newLine.str() << std::endl;
    }
    return newBlockasm.str();
}
