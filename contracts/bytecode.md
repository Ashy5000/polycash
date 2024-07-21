The smart contracts on the blockchain use a special bytecode/machine code specification.

The format outline:

```bytecode
XXXX XXXXXXXX  XXXXXXXX  XXXXXXXX  ... XXXXXXXX  10101010
CMD  PARAMREF0 PARAMREF1 PARAMREF2 ... PARAMREFN PARAMEND
```

| Command      | Code | Param0    | Param1 | Param2 | Param3 | Param4 |
|--------------|------|-----------|--------|--------|--------|--------|
| Exit         | 00   | Status    |
| ExitBfr      | 01   | StatusBfr |
| InitBfr      | 02   | Loc       | ResDst |
| CpyBfr       | 03   | Src       | Dst    | ResDst |
| FreeBfr      | 04   | Loc       | ResDst |
| BfrStat      | 05   | Loc       | StaDst |
| Add          | 06   | A         | B      | Out    | ResDst |
| Sub          | 07   | A         | B      | Out    | ResDst |
| Mul          | 08   | A         | B      | Out    | ResDst |
| Div          | 09   | A         | B      | Out    | ResDst |
| And          | 0A   | A         | B      | Out    | ResDst |
| Or           | 0B   | A         | B      | Out    | ResDst |
| Not          | 0C   | A         | ResDst |
| App          | 0D   | A         | B      | Out    | ResDst |
| Slice        | 0E   | In        | S      | E      | Out    | ResDst |
| Shiftl       | 0F   | In        | Bits   | Out    | ResDst |
| Shiftr       | 10   | In        | Bits   | Out    | ResDst |
| Eq           | 11   | A         | B      | Out    | ResDst |
| Jmp          | 12   | Loc       | ResDst |
| JmpCond      | 13   | Loc       | Cond   | ResDst |
| Stdout       | 14   | Bfr       | ResDst |
| Stderr       | 15   | Bfr       | ResDst |
| SetCnst      | 16   | Bfr       | Val    | ResDst |
| Tx           | 17   | Sndr      | Recvr  | Amount | ResDst |
| GetNthBlk    | 18   | N         | Prop   | Out    | ResDst |
| GetNthTx     | 19   | BlkN      | TxN    | Prop   | Out    | ResDst |
| ChainLen     | 1A   | Out       |
| UpdateState  | 1B   | Location  | Value  | ResDst |
| GetFromState | 1C   | Location  | Dst    | ResDst |

`0x00000000` is a pre-initialized buffer. It is used for capturing global errors.

`Exit`: The `Exit` instruction exits a smart contract with an exit code determined by the `Status` parameter.
`ExitBfr`: The `ExitBfr` instruction 
