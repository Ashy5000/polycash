The smart contracts on the blockchain use a special bytecode/machine code specification.

The format outline:

```bytecode
XXXX XXXXXXXX  XXXXXXXX  XXXXXXXX  ... XXXXXXXX  10101010
CMD  PARAMREF0 PARAMREF1 PARAMREF2 ... PARAMREFN PARAMEND
```

| Command | Code | Param0 | Param1 | Param2 | Param3 |
|---------|------|--------|--------|--------|--------|
| Exit    | 0000 | Status |
| InitBfr | 0001 | Loc    | ResDst |
| CpyBfr  | 0002 | Src    | Dst    | ResDst |
| FreeBfr | 0003 | Loc    | ResDst |
| BfrStat | 0004 | Loc    | StaDst |
| Add     | 0005 | A      | B      | Out    | ResDst |
| Sub     | 0006 | A      | B      | Out    | ResDst |
| Mul     | 0007 | A      | B      | Out    | ResDst |
| And     | 0008 | A      | B      | Out    | ResDst |
| Or      | 0009 | A      | B      | Out    | ResDst |
| Not     | 000A | A      | ResDst |
| App     | 000B | A      | B      | Out    | ResDst |
| Eq      | 000C | A      | B      | Out    | ResDst |
| Jmp     | 000D | Loc    | ResDst |
| JmpCond | 000E | Loc    | Cond   | ResDst |
| Stdout  | 000F | Bfr    | ResDst |
| Stderr  | 0010 | Bfr    | ResDst |
| SetCnst | 0011 | Bfr    | Val    | ResDst |

`0x00000000` is a pre-initialized buffer.