#include "chunk.h"

void InitChunk(Chunk* chunk) {
    DBG("Initing chunk");
    chunk->count = 0;
    chunk->capacity = 0;
    chunk->code = NULL;
    chunk->lines = NULL;
    InitValueArray(&chunk->constants);
}

void WriteChunk(Chunk* chunk, uint8_t byte, int line) {
    if (chunk->count+1 > chunk->capacity) {
        DBG("growing chunk");
        int oldCapacity = chunk->capacity;
        chunk->capacity = GROW_CAPACITY(oldCapacity);

        chunk->code = GROW_ARRAY(uint8_t, chunk->code, oldCapacity, chunk->capacity);
        chunk->lines = GROW_ARRAY(int, chunk->lines, oldCapacity, chunk->capacity);
    }

    DBG("writing byte: '%x' at %p [ %d ]", byte, chunk->code, chunk->count);
    chunk->code[chunk->count] = byte;
    chunk->lines[chunk->count] = line;
    chunk->count++;
}

void FreeChunk(Chunk* chunk) {
    DBG("Freeing chunk code");
    FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
    FREE_ARRAY(int, chunk->lines, chunk->capacity);
    FreeValueArray(&chunk->constants);
    InitChunk(chunk);
}

// XXX: Even though this returns an int, it acts as only a byte because
// it is stored in the chunk as a single byte.
int AddConstant(Chunk *chunk, Value value) {
    DBG("adding constant: %p[%d] = %g", &chunk, chunk->constants.count, value);
    WriteValueArray(&chunk->constants, value);
    return chunk->constants.count-1;
}
