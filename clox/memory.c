#include "memory.h"

#include <stdlib.h>

void *reallocate(void *pointer, size_t oldCapacity, size_t newCapacity) {
    // XXX: why is this required
    if (newCapacity == 0) {
        DBG("Freeing pointer %p of oldCapacity %ld", pointer, oldCapacity);
        free(pointer);
        return NULL;
    }
    DBG("reallocing pointer %p to capacity: %ld", pointer, newCapacity);
    void *ptr = realloc(pointer, newCapacity);
    if (ptr == NULL) {
        DBG("realloc returned NULL");
        exit(1);
    }
    return ptr;
}
