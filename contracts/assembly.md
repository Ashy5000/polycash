For a human-readable version of the bytecode, an assembly language is implemented for smart contracts on the blockchain.

Full command names are used, and `0x` is placed at the beginning of hexadecimal values.

`;` begins a single-line comment

For example:

```
InitBfr 0x00000001 0x00000000 ; Initialize a result buffer
InitBfr 0x00000002 0x00000001 ; Initialize a message buffer
Stdout 0x00000002 0x00000001 ; Send the message buffer's contents to stdout
InitBfr 0x00000003 0x00000001 ; Initialize a status buffer
Exit 0x00000003 0x00000001 ; Exit with code 0
```