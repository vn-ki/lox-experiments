#ifndef clox_common_h
#define clox_common_h

#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

// #define DEBUG 1

#define EC_GREY   "\033[0;37m"

#define EC_RED    "\033[0;31m"
#define EC_L_RED  "\033[1;31m"

#define EC_GREEN    "\033[0;32m"
#define EC_L_GREEN  "\033[1;32m"

#define EC_ORANGE "\033[0;33m"
#define EC_YELLOW "\033[1;33m"

#define EC_PURPLE   "\033[0;35m"
#define EC_L_PURPLE "\033[1;35m"

#define EC_NC     "\033[0m"

#if DEBUG
#include <stdio.h>

#define DBG(M, ...) \
    fprintf(stderr, \
            EC_GREY "%s" EC_NC ":" EC_GREEN "%s" EC_NC ":" EC_ORANGE "%i \t" EC_NC M "\n", \
            __FILE__, __FUNCTION__, __LINE__, ##__VA_ARGS__ \
            ); \
    fflush(stderr);
#else
#define DBG(M, ...)
#endif

#endif
