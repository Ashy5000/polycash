The smart contracts on the blockchain use a special bytecode/machine code specification.

The format outline:

```bytecode
XXXX XXXXXXXX  XXXXXXXX  XXXXXXXX  ... XXXXXXXX  10101010
CMD  PARAMREF0 PARAMREF1 PARAMREF2 ... PARAMREFN PARAMEND
```

| Command | Code | Param0 | Param1 | Param2 |
|---------|------|--------|--------|--------|
| Exit    | 0000 | Status |
| InitBfr |