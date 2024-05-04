The smart contracts on the blockchain use a special bytecode/machine code specification.

The format outline:

```bytecode
XXXX XXXXXXXX  XXXXXXXX  XXXXXXXX  ... XXXXXXXX  10101010
CMD  PARAMREF0 PARAMREF1 PARAMREF2 ... PARAMREFN PARAMEND
```

| Command | Code | Param0 | Param1 | Param2 | Param3 |
| ------- | ---- | ------ | ------ | ------ | ------ | ------ |
| Exit    | 0000 | Status |
| InitBfr | 0001 | Loc    | ResDst |
| CpyBfr  | 0002 | Src    | Dst    | ResDst |
| FreeBfr | 0003 | Loc    | ResDst |
| BfrStat | 0004 | Loc    | StaDst |
| Add     | 0005 | A      | B      | Out    | ResDst |
| Sub     | 0006 | A      | B      | Out    | ResDst |
| Mul     | 0007 | A      | B      | Out    | ResDst |
| Div     | 0008 | A      | B      | Out    | ResDst |
| And     | 0009 | A      | B      | Out    | ResDst |
| Or      | 000A | A      | B      | Out    | ResDst |
| Not     | 000B | A      | ResDst |
| App     | 000C | A      | B      | Out    | ResDst |
| Slice   | 000D | In     | S      | E      | Out    | ResDst |
| Shiftl  | 000E | In     | Bits   | Out    | ResDst |
| Shiftr  | 000F | In     | Bits   | Out    | ResDst |
| Eq      | 0010 | A      | B      | Out    | ResDst |
| Jmp     | 0011 | Loc    | ResDst |
| JmpCond | 0012 | Loc    | Cond   | ResDst |
| Stdout  | 0013 | Bfr    | ResDst |
| Stderr  | 0014 | Bfr    | ResDst |
| SetCnst | 0015 | Bfr    | Val    | ResDst |
| Tx      | 0016 | Sndr   | Recvr  | Amount | ResDst |

`0x00000000` is a pre-initialized buffer.
