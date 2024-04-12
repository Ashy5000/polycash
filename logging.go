// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
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
