#ifndef clox_debug_h
#define clox_debug_h

#include "chunk.h"
#include <stdio.h>

#if DEBUG
#define DISASS_CHUNK(chunk_ptr, name) \
    DBG("Dissembling chunk"); \
    DisassembleChunk(chunk_ptr, name)
#else
#define DISASS_CHUNK(chunk_ptr, name)
#endif

void DisassembleChunk(Chunk *chunk, const char *name);
int DisassembleInstr(Chunk *chunk, int offset);

#endif
