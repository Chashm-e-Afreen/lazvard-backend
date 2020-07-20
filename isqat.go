package main

import (
	"strings"
)

func plurals(word string, dict map[string][]string, lastWord bool) []string {
	temp := word[:len(word)-2]
	weights := []string{}
	if len(dict[temp]) > 0 {
		for i := range weights {
			weights = append(weights, strings.TrimSuffix(weights[i], "0"))
		}
	}
	return weights
}

func lookForIllat(word string, dict map[string][]string, weights []string, lastWord bool) []string {
	if len(weights) == 0 {
		return weights
	}

	exceptionsAlif := []string{}
	exceptionsAlif = append(exceptionsAlif, "کیا", "پا")

	illat := strings.HasSuffix(word, "ا") && !Contains(exceptionsAlif, word) && !lastWord ||
		strings.HasSuffix(word, "ی") && len(word) > 3 && (!lastWord || strings.HasSuffix(word, "ئی")) || //ee can be dropped even if it's at the end
		strings.HasSuffix(word, "ہ") && len(word) > 6 ||
		strings.HasSuffix(word, "و") && len(word) > 3 ||
		strings.HasSuffix(word, "ے") && len(word) > 4 && (!lastWord || strings.HasSuffix(word, "ئے")) ||
		strings.HasSuffix(word, "ؤ") && len(word) > 3 ||
		strings.HasSuffix(word, "ؤں") ||
		((strings.HasSuffix(word, "وں") || strings.HasSuffix(word, "یں")) && len(dict[word[:len(word)-2]]) > 0)

	if illat {
		for i := range weights {
			if strings.HasSuffix(weights[i], "0") {
				weights = append(weights, strings.TrimSuffix(weights[i], string(weights[i][len(weights[i])-1])))
			}
		}

	}

	if strings.HasSuffix(word, "وئے") {
		for i := range weights {
			if strings.HasSuffix(weights[i], "010") {
				weights = append(weights, strings.TrimSuffix(weights[i], "010")+"1")
			}
		}
	}

	return weights
}

func alifeWasl(word string, prevWord string, dict map[string][]string, weights []string, prevWeights []string) ([]string, []string) {

	notValid := strings.HasSuffix(prevWord, "ا") || strings.HasSuffix(prevWord, "ی") || strings.HasSuffix(prevWord, "و") || strings.HasSuffix(prevWord, "ہ") || strings.HasSuffix(prevWord, "ے") || strings.HasSuffix(prevWord, "ۂ") || strings.HasSuffix(prevWord, "ِ")
	if strings.HasPrefix(word, "آ") || strings.HasPrefix(word, "ا") && !notValid {

		for i := range prevWeights { //avoid wasl when the previous word is ke or na
			temp := prevWeights[i][:len(prevWeights[i])-1]
			prevWeights = append(prevWeights, temp) //just have to add adtional prevWeight without the last digit
		}
	}
	return weights, prevWeights
}

func wordVao(word string, prevWord string, dict map[string][]string, weights []string, prevWeights []string) ([]string, []string) {

	constraints := !strings.HasSuffix(prevWord, "ا") && !strings.HasSuffix(prevWord, "ی") && !strings.HasSuffix(prevWord, "ں") && (!strings.HasSuffix(prevWord, "ہ") || strings.HasSuffix(prevWord, "اہ")) && !strings.HasSuffix(prevWord, "و")

	for i := range prevWeights {
		switch {
		case strings.HasSuffix(prevWeights[i], "1") && constraints:
			weights = []string{}
			prevWeights[i] = strings.TrimSuffix(prevWeights[i], "1")
			weights = append(weights, "1")
			weights = append(weights, "10")

		case strings.HasSuffix(prevWeights[i], "0") && constraints:
			weights = []string{}
			prevWeights[i] = strings.TrimSuffix(prevWeights[i], "0")
			weights = append(weights, "1")
			weights = append(weights, "10")

		case !constraints:
			weights = []string{}
			weights = append(weights, "10")
			weights = append(weights, "1")
		}

	}

	return weights, prevWeights
}
