package main

import (
	"regexp"
	"strings"
	"unicode"
)

func splitInput(input string) [][]string {
	input = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(input), "\n") //remove all blank lines
	lines := strings.Split(input, "\n")
	words := [][]string{}

	for i := range lines {
		lines[i] = strings.Join(strings.Fields(strings.TrimSpace(lines[i])), " ")
		words = append(words, strings.Split(lines[i], " "))
	}
	return words
}

func removeNuisances(str string) string {
	chr := " آ آ اب پ ت ٹ ث ج چ ح خ د ڈ ذ ر ڑ ز ژ س ش ص ض ط ظ ع غ ف ق ک گ ل م ن ں و ہ ھ ی ےئ ئے ــِ ةٓ ۂ ۂ ۃ ۂ ؤ ۂ  ً"
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) > 0 || unicode.IsSpace(r) || r == 0x000A || r == 0x100064B { //the last one is for fathatan(do zabar)
			return r
		}
		return -1
	}, str)
}
