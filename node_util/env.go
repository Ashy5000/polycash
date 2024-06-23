// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type NetworkUpgrades struct {
	Guadalajara int `json:"guadalajara"`
	Jinan       int `json:"jinan"`
	Alexandria  int `json:"alexandria"`
}

type Environment struct {
	Network  string          `json:"network"`
	Upgrades NetworkUpgrades `json:"upgrades"`
}

var Env Environment

func LoadEnv() {
	envFile, err := os.Open("env.json")
	if err != nil {
		panic(err)
	}
	defer envFile.Close()
	jsonBytes, _ := ioutil.ReadAll(envFile)
	err = json.Unmarshal(jsonBytes, &Env)
	if err != nil {
		panic(err)
	}
}
