#pragma once

#ifdef __cplusplus
extern "C" {
#endif


typedef unsigned char byte;
typedef unsigned int uint;
#define null 0
typedef unsigned long ulong;
typedef int BOOL;


BOOL cigCalc(byte* pGrappa, byte *pData , int iDataLen , byte *pCig , int *pCigLen);


#ifdef __cplusplus
}
#endif