The smart contracts on the blockchain use a special bytecode/machine code specification.

The format outline:

```bytecode
XXXX XXXXXXXX  XXXXXXXX  XXXXXXXX  ... XXXXXXXX  10101010
CMD  PARAMREF0 PARAMREF1 PARAMREF2 ... PARAMREFN PARAMEND
```

| Command | Code | Param0 | Param1 | Param2 |
|---------|------|--------|--------|--------|
| Exit    | 0000 | Status |
| InitBfr | 0001 | Ptr    | ResDst |
| CpyBfr  | 0002 | Src    | Dst    | ResDst |
| FreeBfr | 0003 | Ptr    | ResDst |
| BfrStat | 0004 | Ptr    | StaDst |
| Add     | 0005 | A      | B      | ResDst |
| Sub     | 0006 | A      | B      | ResDst |
| Mul     | 0007 | A      | B      | ResDst |
| And     | 0008 | A      | B      | ResDst |
| Or      | 0009 | A      | B      | ResDst |
| Not     | 000A | A      | ResDst |