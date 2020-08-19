#include "value.h"

void InitValueArray(ValueArray* value) {
    DBG("Initing value");
    value->count = 0;
    value->capacity = 0;
    value->values = NULL;
}

void WriteValueArray(ValueArray* value, Value v) {
    if (value->count+1 > value->capacity) {
        DBG("growing value");
        int oldCapacity = value->capacity;
        value->capacity = GROW_CAPACITY(oldCapacity);
        value->values = GROW_ARRAY(Value, value->values, oldCapacity, value->capacity);
    }

    DBG("writing byte: '%f' at %p [ %d ]", v, value->values, value->count);
    value->values[value->count] = v;
    value->count++;
}

void FreeValueArray(ValueArray* value) {
    DBG("Freeing value values");
    FREE_ARRAY(Value, value->values, value->capacity);
    InitValueArray(value);
}
