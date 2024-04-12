package main

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func Log(m string, isVerbose bool) {
	m = Blue + "INFO " + Reset + m
	if !isVerbose {
		fmt.Println(m)
		return
	}
	if *verbose {
		fmt.Println(m)
	}
}

func Warn(m string) {
	m = Yellow + "WARNING " + Reset + m
	fmt.Println(m)
}

func Error(m string, fatal bool) {
	m = Red + "ERROR " + Reset + m
	if fatal {
		panic(m)
	} else {
		fmt.Println(m)
	}
}
