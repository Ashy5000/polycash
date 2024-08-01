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

**Exp**

The Exp instruction, represented by the hex instruction code 0x001E, calculates the result of one unsigned 64-bit integer to the power of a second, storing the result in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Exp [A] [B] [Out] [ResDst]

**Mod**

The Exp instruction, represented by the hex instruction code 0x001F, calculates the remainder of dividing two unsigned 64-bit integers in the manner of the Div instruction. If B is equal to zero, the result is zero. The result is stored in the buffer stored at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst

Mod [A] [B] [Out] [ResDst]

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

5.1.3.4 Binary Instructions

**App**

The App instruction, represented by the hex instruction code 0x000C, appends the binary values of the buffers stored at memory locations A and B, storing the result in the buffer at Out. If A, B, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

App [A] [B] [Out] [ResDst]

**Slice**

The Slice instruction, represented by the hex instruction code 0x000D, gets the contents of the registers at locations S and E, converts them to unsigned 64-bit integers, and uses them as the start and end indexes of a range of bits from In. The contents of this range is copied to Out. If In, S, E, or Out do not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Slice [In] [S] [E] [Out] [ResDst]

**Shiftl**

The Shiftl instruction, represented by the hex instruction code 0x000E, applies the left binary shift operation to the buffer located at In by a number of bits specified by an unsigned 64-bit integer located in the buffer at memory location Bits. The result is stored in Out. If In, Bits, or Out o not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Shiftl [In] [Bits] [Out] [ResDst]

**Shiftr**

The Shiftr instruction, represented by the hex instruction code 0x000F, applies the right binary shift operation to the buffer located at In by a number of bits specified by an unsigned 64-bit integer located in the buffer at memory location Bits. The result is stored in Out. If In, Bits, or Out o not hold memory locations with an initialized buffer, a local error is thrown into ResDst.

Shiftr [In] [Bits] [Out] [ResDst]

**Eq**

The Eq instruction, represented by the hex instruction code 0x0010, checks if the buffers at memory locations A and B are equal. If they are, the value 0x0000000000000001 is stored in the buffer Out. Otherwise, 0x0000000000000000 is stored. if A, B, or Out do not exist, a local error is thrown into ResDst.

Eq [A] [B] [Out] [ResDst]

**5.1.3.5** Control Flow Instructions

**Jmp**

The Jmp instruction, represented by the hex instruction code 0x0011, jumps to the instruction specified by the line number stored in the buffer at Loc. If Loc does not exist, a local error is thrown into ResDst.

Jmp [Loc] [ResDst]

**JmpCond**

The JmpCond instruction, represented by the hex instruction code 0x0012, jumps to the instruction specified by the line number Loc if the buffer at Cond stores a value other than 0x0000000000000000. If Loc or Cond do not exist, a local error is thrown into ResDst.

JmpCond [Loc] [Cond] [ResDst]

**Call**

The Call function, represented by the hex instruction code 0x0021, jumps to the line Loc. If Loc does not exist, a local error is thrown into ResDst.

Call [Loc] [ResDst]

**Ret**

The Ret instruction, represented by the hex instruction code 0x0022, jumps to the line where the Call instruction was last run.

Ret

**5.1.3.5** Contract instructions

**Exit**

The Exit instruction, represented by the hex instruction code 0x0000, returns from the smart contract with exit code Status.

Exit [Status]

**ExitBfr**

The ExitBfr instruction, represented by the hex instruction code 0x0001, returns from the smart contract with the exit code stored in the buffer at StatusBfr. If StatusBfr does not exist, a local error is thrown into ResDst.

ExitBfr [StatusBfr] [ResDst]

**5.1.3.5** Logging instructions

**Stdout**

The Stdout instruction, represented by the hex instruction code 0x0014, prints the data in the buffer at the location Bfr in a debugging format. If Bfr does not exist, a local error is thrown into ResDst.

Stdout [Bfr] [ResDst]

**Stderr**

The Stderr instruction, represented by the hex instruction code 0x0015, prints the data in the buffer at the location Bfr in a debugging format into STDERR. If Bfr does not exist, a local error is thrown into ResDst.

Stdout [Bfr] [ResDst]

**PrintStr**

The PrintStr instruction, represented by the hex instruction code 0x0023, prints the ASCII string stored in the buffer at StrBfr. If StrBfr does not exist, a local error is thrown into ResDst.

PrintStr [StrBfr] [ResDst]
