#ifndef clox_memory_h
#define clox_memory_h

#include "common.h"

#define GROW_CAPACITY(c) \
    c < 8 ? 8 : c * 2
#define GROW_ARRAY(type, ptr, oldCount, newCapacity) \
    (type *)reallocate(ptr, sizeof(type) * oldCount, sizeof(type) * newCapacity)
#define FREE_ARRAY(type, ptr, capacity) \
    reallocate(ptr, sizeof(type) * capacity, 0)

void *reallocate(void *pointer, size_t oldCapacity, size_t newCapacity);


#endif
