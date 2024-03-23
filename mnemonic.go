package main

import (
	"crypto/dsa"
	"fmt"
	"math/big"
	"strings"
)

var words = map[string]string{
	"00": "abandon",
	"01": "ability",
	"02": "able",
	"03": "about",
	"04": "above",
	"05": "absent",
	"06": "absorb",
	"07": "abstract",
	"08": "absurd",
	"09": "abuse",
	"10": "access",
	"11": "accident",
	"12": "account",
	"13": "accuse",
	"14": "achieve",
	"15": "acid",
	"16": "acoustic",
	"17": "acquire",
	"18": "across",
	"19": "act",
	"20": "action",
	"21": "actor",
	"22": "actress",
	"23": "actual",
	"24": "adapt",
	"25": "add",
	"26": "addict",
	"27": "address",
	"28": "adjust",
	"29": "admit",
	"30": "adult",
	"31": "advance",
	"32": "advice",
	"33": "aerobic",
	"34": "affair",
	"35": "afford",
	"36": "afraid",
	"37": "again",
	"38": "age",
	"39": "agent",
	"40": "agree",
	"41": "ahead",
	"42": "aim",
	"43": "air",
	"44": "airport",
	"45": "aisle",
	"46": "alarm",
	"47": "album",
	"48": "alcohol",
	"49": "alert",
	"50": "alien",
	"51": "all",
	"52": "alley",
	"53": "allow",
	"54": "almost",
	"55": "alone",
	"56": "alpha",
	"57": "already",
	"58": "also",
	"59": "alter",
	"60": "always",
	"61": "amateur",
	"62": "amazing",
	"63": "among",
	"64": "amount",
	"65": "amused",
	"66": "analyst",
	"67": "anchor",
	"68": "ancient",
	"69": "anger",
	"70": "angle",
	"71": "angry",
	"72": "animal",
	"73": "ankle",
	"74": "announce",
	"75": "annual",
	"76": "another",
	"77": "answer",
	"78": "antenna",
	"79": "antique",
	"80": "anxiety",
	"81": "any",
	"82": "apart",
	"83": "apology",
	"84": "appear",
	"85": "apple",
	"86": "approve",
	"87": "april",
	"88": "arch",
	"89": "arctic",
	"90": "area",
	"91": "arena",
	"92": "argue",
	"93": "arm",
	"94": "armed",
	"95": "armor",
	"96": "army",
	"97": "around",
	"98": "arrange",
	"99": "arrest",
	// Words with placeholders
	"0!": "buddy",
	"1!": "basic",
	"2!": "bizarre",
	"3!": "bless",
	"4!": "boring",
	"5!": "brave",
	"6!": "breeze",
	"7!": "bribe",
	"8!": "broken",
	"9!": "banana",
}

func GetMnemonic(key dsa.PrivateKey) string {
	// Get each group of 2 digits from the key and convert it to a word.
	// Convert the key to a string.
	keyString := key.X.String()
	// Split the string into groups of 2 digits.
	var groups []string
	if len(keyString)%2 != 0 {
		keyString = keyString + "!"
	}
	for i := 0; i < len(keyString); i += 2 {
		groups = append(groups, keyString[i:i+2])
	}
	fmt.Println(groups)
	// Convert each group of 2 digits to a word.
	var mnemonic string

	for _, group := range groups {
		mnemonic += words[group] + " "
	}

	fmt.Println(mnemonic)
	return mnemonic
}

func RestoreMnemonic(mnemonic string) dsa.PrivateKey {
	// Split the mnemonic into words.
	mnemonicWords := strings.Split(mnemonic, " ")
	// Convert each word to a group of 2 digits.
	groups := make([]string, len(words))
	for _, word := range mnemonicWords {
		for key, value := range words {
			if value == word {
				groups = append(groups, key)
			}
		}
	}
	// Join the groups of 2 digits together.
	keyString := strings.Join(groups, "")
	// Remove the last character if it is a "!".
	if keyString[len(keyString)-1] == '!' {
		keyString = keyString[:len(keyString)-1]
	}
	// Convert the string to a big.Int.
	key := dsa.PrivateKey{
		PublicKey: dsa.PublicKey{},
		X:         big.NewInt(0),
	}
	key.X.SetString(keyString, 10)
	return key
}
