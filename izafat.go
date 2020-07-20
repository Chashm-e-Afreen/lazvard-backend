package main

import (
	"strings"
)

func lookForIzafat(word string, dict map[string][]string, weights []string) []string {
	kasra := "ِ"
	izafatYa := "ئے"
	izafatHamza := "ۂ"
	izafat := strings.HasSuffix(word, kasra) ||
		strings.HasSuffix(word, izafatYa) ||
		strings.HasSuffix(word, izafatHamza)

	if izafat && !strings.HasSuffix(word, "ں") {
		switch {
		case strings.HasSuffix(word, kasra):
			weights = foundIzafat(word, dict, kasra)
		case strings.HasSuffix(word, "ائے") || strings.HasSuffix(word, "وئے"):
			weights = foundIzafat(word, dict, izafatYa)
		case strings.HasSuffix(word, izafatHamza):
			weights = foundIzafat(word, dict, izafatHamza)
		}

	}
	return weights
}

func foundIzafat(word string, dict map[string][]string, suffix string) []string {
	rootWord := strings.TrimSuffix(word, suffix)
	if suffix == "ۂ" || suffix == "ۂ" {
		rootWord += "ہ"
	}
	temp := dict[rootWord]
	weights := []string{}
	for i := range temp {
		switch {
		case strings.HasSuffix(temp[i], "1"):
			weights = append(weights, temp[i])
			weights = append(weights, temp[i]+"0")
		default:
			if suffix == "ئے" {
				local := temp[i]
				weights = append(weights, local+"10")
				weights = append(weights, local+"1")
			} else {
				local := strings.TrimSuffix(temp[i], "0")
				weights = append(weights, local+"10")
				weights = append(weights, local+"1")
			}
		}
	}
	return weights
}
