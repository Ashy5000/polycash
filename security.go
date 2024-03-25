package main

type SecurityLevel struct {
	Level                  int
	MinimumDifficulty      uint64
	InitialBlockDifficulty uint64
	BlocksBeforeSpendable  int
}

var securityLevels = []SecurityLevel{
	{0, 100000, 100000, 3},
	{1, 200000, 250000, 5},
	{2, 500000, 500000, 7},
}

func ApplySecurityLevel(level int) {
	for _, securityLevel := range securityLevels {
		if securityLevel.Level == level {
			initialBlockDifficulty = securityLevel.InitialBlockDifficulty
			minimumBlockDifficulty = securityLevel.MinimumDifficulty
			blocksBeforeSpendable = securityLevel.BlocksBeforeSpendable
		}
	}
}
