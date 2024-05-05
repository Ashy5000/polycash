**Polycash VM Instruction Set**

This document is an extended subsection removed from the original Polycash whitepaper.

**5.1.3 Instruction Set**

Details for the instruction set for the Polycash-specific virtual machine is covered in the following subsections.

**5.1.3.1 Buffer Instructions**

**InitBfr**

The InitBfr instruction, represented by the hex instruction code 0x0001, initializes a buffer at an address specified by the provided Loc buffer. If there is already a buffer at that address, no buffer is created and a local error is thrown into the provided ResDst buffer.

InitBfr [Loc] [ResDst]

**FreeBfr**

The FreeBfr instruction, represented by the hex instruction code 0x0002;, frees a buffer at an address specified by the provided Loc buffer. If there is no buffer at that address, no buffer is freed and a local error is thrown into the provided ResDst buffer.

FreeBfr [Loc] [ResDst]

**CpyBfr**

The CpyBfr instruction, represented by the hex instruction code 0x0003;, copies a buffer at an address specified by the provided Src buffer to that of the Dst buffer. If either buffer does not exist, no buffer is copied and a local error is thrown into the provided ResDst buffer.

CpyBfr [Src] [Dst] [ResDst]

**SetCnst**

The SetCnst instruction, represented by the hex instruction code 0x0015, copies a constant value, Val, into a buffer located at the address stored in the Bfr buffer. If the buffer does not exist, a local error is thrown into ResDst. Please note that “constant” refers to the fact that the value is determined at compile time. A buffer with a value set by SetCnst is mutable.

SetCnst [Bfr] [Val] [ResDst]

**BfrStat**

The BfrStat instruction, represented by the hex instruction code 0x0004, gets the initialization status at the address Loc. If there is an initialized buffer at that location, an unsigned 8-bit integer representation of the number 1 is stored in the buffer at location Dst. Otherwise, a 0 is stored. If there is no buffer located at Dst, a local error is thrown into ResDst.

BfrStat [Loc] [Dst] [ResDst]

5.1.3.2 Math Instructions

**Add**

The Add instruction, represented by the hex instruction code 0x0005, adds two unsigned 64-bit integers stored in the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Add [A] [B] [Out] [ResDst]

**Sub**

The Sub instruction, represented by the hex instruction code 0x0006, subtracts two unsigned 64-bit integers stored at the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Sub [A] [B] [Out] [ResDst]

**Mul**

The Mul instruction, represented by the hex instruction code 0x0007, multiplies two unsigned 64-bit integers stored at the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Mul [A] [B] [Out] [ResDst]

**Div**

The Div instruction, represented by the hex instruction code 0x0008, divides two unsigned 64-bit integers stored at the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Div [A] [B] [Out] [ResDst]

5.1.3.3 Logical Instructions

**And**

The And instruction, represented by the hex instruction code 0x0009, performs a bitwise and operation on two 64-bit values stored at the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Div [A] [B] [Out] [ResDst]

**Or**

The Or instruction, represented by the hex instruction code 0x000A, performs a bitwise or operation on two 64-bit values stored at the buffers at the memory locations A and B, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Or [A] [B] [Out] [ResDst]

**Not**

The Not instruction, represented by the hex instruction code 0x000B, stores 0x0000000000000001 in Out if A is equal to 0x0000000000000000, and 0x0000000000000000 otherwise. If A or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Not [A] [Out] [ResDst]
