package main

import (
	"strings"
)

func checkRoot(word string, dict map[string][]string, lastWord bool) []string {
	weights := []string{}
	switch {
	case strings.HasSuffix(word, "ے"):
		weights = foundBariYe(word, dict, weights)
	case strings.HasSuffix(word, "ں"):
		weights = foundNoonGhunna(word, dict, weights, lastWord)
	case strings.HasSuffix(word, "ی"):
		weights = foundYe(word, dict, weights)
	case strings.HasSuffix(word, "و"):
		weights = foundVao(word, dict, weights)
	case strings.HasSuffix(word, "ن"):
		weights = foundNoon(word, dict, weights)
	default:
		return weights
	}
	return weights
}

func foundBariYe(word string, dict map[string][]string, weights []string) []string {

	switch {
	case strings.HasSuffix(word, "ئے"):
		ey := "ئے"
		switch {
		case len(word) < 4:
			return weights
		default:
			rootWord := strings.TrimSuffix(word, ey)
			temp := dict[rootWord]

			weights := []string{}
			if len(temp) != 0 {
				for i := range temp {
					local := temp[i]
					weights = append(weights, local+"1")
					weights = append(weights, local+"10")
				}
				return weights
			}

		}
	case !strings.HasSuffix(word, "ئے"):
		rootWord := strings.TrimSuffix(word, "ے")
		temp := dict[rootWord]
		weights = []string{}
		for i := range temp {
			local := temp[i]
			if strings.HasSuffix(temp[i], "0") {

				local = strings.TrimSuffix(temp[i], "0")
				local = temp[i] + "1"
			}
			weights = append(weights, local+"0")
		}
		if len(weights) == 0 {
			rootWord := strings.TrimSuffix(word, "ے")
			rootWord += "ا"
			temp := dict[rootWord]
			weights = []string{}
			for i := range temp {
				local := temp[i]
				weights = append(weights, local)
			}
			if len(weights) == 0 {
				rootWord := strings.TrimSuffix(word, "ا")
				rootWord += "ہ"
				temp := dict[rootWord]
				weights = []string{}
				for i := range temp {
					local := temp[i]
					weights = append(weights, local)
				}
			}
		}
	default:

		return weights
	}
	return weights
}

func foundNoonGhunna(word string, dict map[string][]string, weights []string, lastWord bool) []string {

	switch {
	case strings.HasSuffix(word, "یاں"):
		weights = foundYaan(word, dict, weights, lastWord)
	case strings.HasSuffix(word, "اں"):
		weights = foundAan(word, dict, weights)
	case strings.HasSuffix(word, "یں") && len(word) > 6: //the izafat method just happens to work with these suffixes
		weights = foundEnOun(word, dict, "یں", lastWord)
	case strings.HasSuffix(word, "وں") && len(word) > 6:
		weights = foundOn(word, dict, lastWord)
	case strings.HasSuffix(word, "ؤں") && len(word) > 6:
		weights = foundEnOun(word, dict, "ؤں", lastWord)
	default:
		noonGhunna := "ں"
		rootWord := strings.TrimSuffix(word, noonGhunna) + "ن"
		temp := dict[rootWord]
		weights = []string{}
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

func foundYaan(word string, dict map[string][]string, weights []string, lastWord bool) []string {
	yaan := "یاں"
	switch {
	case len(word) < 4:
		return weights
	default:
		rootWord := strings.TrimSuffix(word, yaan)
		temp := dict[rootWord]
		if len(temp) != 0 {
			for i := range temp {
				local := temp[i]
				if !lastWord {
					local += "1"
					weights = append(weights, local)
					local += "0"
					weights = append(weights, local)
				} else {
					local += "10"
					weights = append(weights, local)
				}
			}
			return weights
		}
		return weights

	}
}

func foundAan(word string, dict map[string][]string, weights []string) []string {
	aan := "اں"
	switch {
	case len(word) < 4:
		return weights

	case strings.HasSuffix(word, "گاں"):
		switch {
		case len(word) < 4:
			return weights
		default:
			rootWord := strings.TrimSuffix(word, "گاں")
			rootWord += "ہ"
			temp := dict[rootWord]
			if len(temp) != 0 {
				for i := range temp {
					local := temp[i]
					local = strings.TrimSuffix(local, "0")
					weights = append(weights, local+"10")
				}
				return weights
			}
			return weights

		}

	default:
		rootWord := strings.TrimSuffix(word, aan)
		temp := dict[rootWord]
		if len(temp) != 0 {
			for i := range temp { //remove last character and add 10
				local := temp[i]
				if strings.HasSuffix(local, "0") {
					local = strings.TrimSuffix(local, string(local[len(local)-1]))
					weights = append(weights, local+"10")
				} else {
					weights = append(weights, local+"0")
				}
			}

		} else {
			noonGhunna := "ں"
			rootWord := strings.TrimSuffix(word, noonGhunna) + "ن"
			temp := dict[rootWord]
			weights = []string{}
			if len(temp) != 0 {
				for i := range temp {
					local := temp[i]
					weights = append(weights, strings.TrimSuffix(local, "1"))
				}
				return weights
			}
		}

	}
	return weights
}

func foundYe(word string, dict map[string][]string, weights []string) []string {
	switch {
	case len(word) < 3:
		return weights
	case strings.HasSuffix(word, "گی"):
		rootWord := strings.TrimSuffix(word, "گی")
		rootWord += "ہ"
		temp := dict[rootWord]
		if len(temp) != 0 {
			for i := range temp {
				local := temp[i]
				local = strings.TrimSuffix(local, "0")
				weights = append(weights, local)
			}

		} else {
			return weights
		}
	default:
		rootWord := strings.TrimSuffix(word, "ی")
		temp := dict[rootWord]
		if len(temp) != 0 {
			for i := range temp {
				local := temp[i]
				local = strings.TrimSuffix(local, string(local[len(local)-1]))
				weights = append(weights, local+"10")

			}

		} else {
			noon := "ن"
			derivative := strings.TrimSuffix(word, noon) + "ں"
			temp := dict[derivative]
			if len(temp) != 0 {
				for i := range temp {
					local := temp[i]
					weights = append(weights, local+"1")
				}
				return weights
			}
			return weights

		}

	}
	return weights
}

func foundVao(word string, dict map[string][]string, weights []string) []string {
	vao := "و"
	switch {
	case len(word) < 3:
		return weights
	default:
		rootWord := strings.TrimSuffix(word, vao)
		temp := dict[rootWord]
		if len(temp) != 0 {
			for i := range temp {
				local := temp[i]
				local = strings.TrimSuffix(local, string(local[len(local)-1]))
				weights = append(weights, local+"10")
			}
			return weights
		}
		return weights
	}
}

func foundEnOun(word string, dict map[string][]string, suffix string, lastWord bool) []string {
	switch {
	case strings.HasSuffix(word, "ئیں") || strings.HasSuffix(word, "ؤں"):

		rootWord := strings.TrimSuffix(word, suffix)
		if strings.HasSuffix(rootWord, "ئ") {
			rootWord = strings.TrimSuffix(rootWord, "ئ")
		}
		if strings.HasSuffix(word, "ئ") {
			rootWord = strings.TrimSuffix(word, "ئ")
		}
		temp := dict[rootWord]
		weights := []string{}
		for i := range temp {
			local := temp[i]
			if !lastWord {
				weights = append(weights, local+"1")
			}
			weights = append(weights, local+"10")
		}
		return weights
	default:
		weights := foundIzafat(word, dict, "یں")
		return weights
	}
}

func foundNoon(word string, dict map[string][]string, weights []string) []string {

	noon := "ن"
	derivative := strings.TrimSuffix(word, noon) + "ں"
	temp := dict[derivative]
	if len(temp) != 0 {
		for i := range temp {
			local := temp[i]
			weights = append(weights, local+"1")
		}
		return weights
	}
	return weights
}

func foundOn(word string, dict map[string][]string, lastWord bool) []string {
	rootWord := strings.TrimSuffix(word, "وں")
	temp := dict[rootWord]
	weights := []string{}

	if len(temp) > 0 {

		for i := range temp {
			local := temp[i]
			switch {
			case strings.HasSuffix(local, "1"):
				weights = append(weights, local)
				if !lastWord {
					weights = append(weights, local+"0")
				}
			default:
				local = strings.TrimSuffix(local, "0")
				if !lastWord {
					weights = append(weights, local+"1")
				}
				weights = append(weights, local+"10")
			}
		}
	} else {
		rootWord = rootWord + "ہ"
		if len(dict[rootWord]) > 0 {
			temp1 := dict[rootWord]
			for i := range temp1 {
				weights = append(weights, temp1[i])
			}
		} else {
			rootWord = strings.TrimSuffix(rootWord, "ہ") + "ا"
			if len(dict[rootWord]) > 0 {
				temp1 := dict[rootWord]
				for i := range temp1 {
					weights = append(weights, temp1[i])
				}
			}
		}
	}
	return weights
}
