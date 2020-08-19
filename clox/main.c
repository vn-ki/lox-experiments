#include <stdio.h>

#include "common.h"
#include "chunk.h"
#include "debug.h"

int main(int argc, const char* argv[]) {
    Chunk chunk;

    InitChunk(&chunk);

    int constant = AddConstant(&chunk, 2);
    WriteChunk(&chunk, OP_CONSTANT, 1);
    WriteChunk(&chunk, constant, 1);
    WriteChunk(&chunk, OP_RETURN, 1);
    DISASS_CHUNK(&chunk, "test");
    FreeChunk(&chunk);
    return 0;
}
