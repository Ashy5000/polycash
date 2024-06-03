package analysis

import (
	. "cryptocurrency/node_util"
	"time"
)

func GetTPS(duration time.Duration) float64 {
	now := time.Now()
	txCount := 0
	for i := len(Blockchain) - 1; i >= 0; i-- {
		block := Blockchain[i]
		if now.Sub(block.Timestamp) > duration {
			break
		}
		txCount += len(block.Transactions)
	}
	return float64(txCount) / duration.Seconds()
}
