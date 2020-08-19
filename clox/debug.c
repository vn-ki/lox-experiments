#include "debug.h"


void DisassembleChunk(Chunk *chunk, const char *name) {
    printf(EC_PURPLE "=====%s=====\n" EC_NC, name);

    for (int offset=0; offset<chunk->count;) {
        offset = DisassembleInstr(chunk, offset);
    }
}

int simpleInstruction(char *name, int offset) {
    printf("%s\n", name);
    return offset+1;
}

void printValue(Value val) {
    printf("%g", val);
}

int constantInstr(char *name, Chunk *chunk, int offset) {
    uint8_t constant = chunk->code[offset+1];
    printf("%-16s %4d '", name, constant);
    printValue(chunk->constants.values[constant]);
    printf("'\n");
    return offset + 2;
}

int DisassembleInstr(Chunk *chunk, int offset) {
    printf("%04d ", offset);
    if (offset > 0 && chunk->lines[offset] == chunk->lines[offset - 1]) {
        printf("   | ");
    } else {
        printf("%4d ", chunk->lines[offset]);
    }

    uint8_t instr = chunk->code[offset];
    switch (instr) {
        case OP_RETURN:
            return simpleInstruction("OP_RETURN", offset);
        case OP_CONSTANT:
            return constantInstr("OP_CONSTANT", chunk, offset);
        default:
            printf("Unknown opcode '%d'\n", instr);
            return offset+1;
    }
}

