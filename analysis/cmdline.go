package analysis

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var analysisCommands = map[string]func([]string){
	"GetTPS":          GetTPSCmd,
	"GetTokensMinted": GetTokensMintedCmd,
}

func GetTPSCmd(fields []string) {
	duration := fields[1]
	durationInt, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		panic(err)
	}
	tps := GetTPS(time.Duration(durationInt) * time.Second)
	fmt.Println(tps)
}

func GetTokensMintedCmd(fields []string) {
	tokensMinted := GetNumTokensMinted()
	fmt.Println(tokensMinted)
}

func RunAnalysisCmd(input string) {
	cmds := strings.Split(input, ";")
	for _, cmd := range cmds {
		fields := strings.Split(cmd, " ")
		action := fields[0]
		switch action {
		case "exit":
			os.Exit(0)
		case "":
			return
		default:
			fn, ok := analysisCommands[action]
			if ok {
				fn(fields)
			} else {
				fmt.Println("Invalid command")
			}
		}
	}
}

func StartAnalysisCmdline() {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Print("analysis> ")
		cmd, _ := inputReader.ReadString('\n')
		cmd = cmd[:len(cmd)-1]
		RunAnalysisCmd(cmd)
	}
}
