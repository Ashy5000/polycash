#include "Compiler.h"

#include <algorithm>

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
    std::string generatedBlockasm = generator.GenerateBlockasm(cm);
    int nextAllocatedLocation = generator.GetNextAllocatedLocation();
    std::string controlSegment = cm.compile(nextAllocatedLocation);
    controlSegmentSize = std::count( controlSegment.begin(), controlSegment.end(), '\n' ) ;
    controlSegmentSize += !controlSegment.empty() && controlSegment.back() != '\n';
    std::string controlClose = cm.close(nextAllocatedLocation);
    blockasm = controlSegment + generatedBlockasm + controlClose;
}

void Compiler::Link() {
    // NOTE: The names Linker and LabelManager are somewhat misleading.
    // The Linker is used during the GenerateBlockasm step to inject function calls and source code.
    // The LabelManager is used during the Link step to resolve uses of the Jmp, JmpCond, and Call instructions.
    // Both the Linker and the LabelManager perform linking, just at different times.
    auto lm = LabelManager(blockasm);
    blockasm = lm.SkipLibs(blockasm);
    blockasm = lm.ReplacePreLabels(blockasm);
    blockasm = lm.ReplaceLabels(blockasm);
    blockasm = lm.OffsetCalls(blockasm, controlSegmentSize);
}

std::string Compiler::Compile(std::string filename_p) {
    filename = filename_p;
    LoadContents();
    ParseTokens();
    GenerateBlockasm();
    Link();
    return blockasm;
}
