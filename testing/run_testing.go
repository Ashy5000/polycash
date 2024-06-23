package testing

import (
	. "cryptocurrency/node_util"
	"fmt"
	"time"
)

func StartTest() {
	StartNetwork()
	tpsIn := 1
	secs := 5
	SendTxs(int64(tpsIn), int64(secs))
	start := time.Now()
	initialtxs := 0
	for _, block := range Blockchain {
		initialtxs += len(block.Transactions)
	}
	for {
		SyncBlockchain(-1)
		txs := 0
		for _, block := range Blockchain {
			txs += len(block.Transactions)
		}
		if txs-initialtxs >= tpsIn*secs {
			duration := time.Since(start)
			fmt.Printf("TPS: %f\n", float64(tpsIn*secs)/duration.Seconds())
		}
		time.Sleep(time.Second * 5)
	}
}
