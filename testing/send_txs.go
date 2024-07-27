package testing

import (
	. "cryptocurrency/node_util"
	"fmt"
	"time"
)

func SendTxs(rate int64, seconds int64) {
	delay := time.Second / time.Duration(rate)
	for i := int64(0); i < seconds*rate; i++ {
		Send("YWJj", "0", []byte(fmt.Sprint(i)))
		time.Sleep(delay)
	}
}
