#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include "memory.h"
#include "value.h"

typedef enum {
    OP_RETURN,
    OP_CONSTANT
} OpCode;

typedef struct {
    int count;
    int capacity;
    uint8_t* code;
    int* lines;
    ValueArray constants;
} Chunk;

void InitChunk(Chunk *chunk);
void WriteChunk(Chunk *chunk, uint8_t byte, int line);
void FreeChunk(Chunk *chunk);
int AddConstant(Chunk *chunk, Value);

#endif
