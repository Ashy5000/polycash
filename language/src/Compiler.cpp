#include "Compiler.h"
#include "Parser.h"
#include "BlockasmGenerator.h"
#include "LabelManager.h"
#include <sstream>
#include <fstream>

void Compiler::LoadContents() {
    std::stringstream contents_stream;
    {
        std::fstream input(filename, std::ios::in);
        contents_stream << input.rdbuf();
    }
    contents = contents_stream.str();
}

void Compiler::ParseTokens() {
    tokens = Parser::parse_tokens(contents);
}

void Compiler::GenerateBlockasm() {
    auto generator = BlockasmGenerator(tokens, 0x00001000, {}, true);
    auto cm = ControlModule();
    int nextAllocatedLocation = generator.GetNextAllocatedLocation();
    std::string generatedBlockasm = generator.GenerateBlockasm(cm);
    std::string controlSegment = cm.compile(nextAllocatedLocation);
    blockasm = controlSegment + generatedBlockasm;
}

void Compiler::Link() {
    auto lm = LabelManager(blockasm);
    blockasm = lm.ReplacePreLabels(blockasm);
    blockasm = lm.ReplaceLabels(blockasm);
}

std::string Compiler::Compile(std::string filename_p) {
    filename = filename_p;
    LoadContents();
    ParseTokens();
    GenerateBlockasm();
    Link();
    return blockasm;
}
