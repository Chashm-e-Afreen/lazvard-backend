package main

import "strings"

func fetchWeights(word string, dict map[string][]string, lastWord bool) []string {

	izafat := strings.HasSuffix(word, "ِ")
	word = strings.ReplaceAll(word, "ِ", "")

	if izafat {
		word += "ِ"
	}
	weights := copyString1d(dict[word])
	if lastWord {
		for i := range weights {
			if i < len(weights) {
				if weights[i] == "1" {
					weights[i], weights[len(weights)-1] = weights[len(weights)-1], weights[i]
					weights = weights[:len(weights)-1]
				}
			}
		}
	}

	if len(word) > 1 {
		if !strings.HasSuffix(word, "ئے") || len(word) > 8 { //to avoid aaye, jaye being interpreted as izafats
			weights = lookForIzafat(word, dict, weights)
		}
	}
	if len(weights) == 0 {
		weights = checkRoot(word, dict, lastWord)
	}

	if len(weights) != 0 && len(word) > 3 { //lastWord of the line
		if !lastWord {
			weights = lookForIllat(word, dict, weights, lastWord)
		}
	}

	if len(weights) == 0 {
		if strings.HasSuffix(word, "وں") || strings.HasSuffix(word, "یں") {
			weights = plurals(word, dict, lastWord)
		}
	}

	if len(weights) == 0 {
		noonGhunna := "ں"
		rootWord := strings.TrimSuffix(word, noonGhunna) + "ن"
		temp := dict[rootWord]
		weights := []string{}
		if len(temp) != 0 {
			for i := range temp {
				local := temp[i]
				weights = append(weights, strings.TrimSuffix(local, "1"))
			}
			return weights
		}
	}

	return weights
}

func lineScansion(words [][]string, dict map[string][]string) [][][]string {

	lineWeights := make([][][]string, len(words))

	for i := range words {

		for j := range words[i] {

			lastWord := (j == len(words[i])-1)

			if j > 0 && (strings.HasPrefix(words[i][j], "آ") || strings.HasPrefix(words[i][j], "ا") || words[i][j] == "و") {

				lineWeights[i] = append(lineWeights[i], fetchWeights(words[i][j], dict, lastWord))

				vaslConstraints := !strings.HasSuffix(words[i][j-1], "ا") && !strings.HasSuffix(words[i][j-1], "ی") && !strings.HasSuffix(words[i][j-1], "ں")

				if vaslConstraints {

					lineWeights[i][j], lineWeights[i][j-1] = alifeWasl(words[i][j], words[i][j-1], dict, lineWeights[i][j], lineWeights[i][j-1])
				}

				if words[i][j] == "و" {
					lineWeights[i][j], lineWeights[i][j-1] = wordVao(words[i][j], words[i][j-1], dict, lineWeights[i][j], lineWeights[i][j-1])
				}

			} else {
				lineWeights[i] = append(lineWeights[i], fetchWeights(words[i][j], dict, lastWord))
			}

		}
	}
	return lineWeights
}
