#ifndef clox_value_h
#define clox_value_h

#include "common.h"
#include "memory.h"

typedef double Value;

typedef struct {
    int count;
    int capacity;
    Value *values;
} ValueArray;

void InitValueArray(ValueArray *);
void WriteValueArray(ValueArray *, Value);
void FreeValueArray(ValueArray *);

#endif
