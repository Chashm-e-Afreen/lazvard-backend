package main

import (
	// "strconv"

	"strings"

	levenshtein "github.com/agnivade/levenshtein"
	levenshtein1 "github.com/texttheater/golang-levenshtein/levenshtein"
)

func genKey(combination []string) (string, []string) {

	combinationKey := strings.Join(combination[:], "")
	return combinationKey, combination
}

func meterMatch(combinationKey string, bestMatch []string, sliceCombination []string, meters map[string]string, meterNames map[string]string) (int, string, string, string, bool, int) {

	tasbeeghOazala := 0
	combinationKey, tasbeeghOazala = checkTasbeeghAzala(combinationKey)
	closestMeter := meters[combinationKey]
	closestMeterName := meterNames[combinationKey]
	dist := 0
	var perfectMatch bool
	if len(closestMeter) != 0 && (Contains(bestMatch, combinationKey) || len(bestMatch) == 0) {
		perfectMatch = true
		switch {
		case tasbeeghOazala == 1:
			closestMeter = strings.TrimSuffix(closestMeter, "ن")
			closestMeter += "ان"
			closestMeterName += " مسبغ"
		case tasbeeghOazala == 2:
			closestMeter = strings.TrimSuffix(closestMeter, "ن")
			closestMeter += "ان"
			closestMeterName += "مذال"
		}
		return dist, closestMeter, closestMeterName, combinationKey, perfectMatch, tasbeeghOazala
	}
	perfectMatch = false

	closestMeterKey := ""

	switch {

	case len(bestMatch) != 0:
		for i := range bestMatch {
			closestMeterKey = bestMatch[i]
			temp := levenshtein.ComputeDistance(combinationKey, bestMatch[i])
			if temp < dist || i == 0 {
				dist = temp
				closestMeterKey = bestMatch[i]
			}
		}

	default:
		for i := range meterList {
			temp := levenshtein.ComputeDistance(combinationKey, meterList[i])
			if temp < dist || i == 0 {
				dist = temp
				closestMeterKey = meterList[i]
			}
		}

	}

	closestMeter = meters[closestMeterKey]
	closestMeterName = meterNames[closestMeterKey]
	dist = levenshtein.ComputeDistance(combinationKey, closestMeterKey)

	if !perfectMatch && checkMuqatta(closestMeterKey) && dist < 2 { //tasbeeghOazala for muqatta bahoor
		editScript := levenshtein1.EditScriptForStrings([]rune(closestMeterKey), []rune(combinationKey), levenshtein1.DefaultOptions)

		if editScript[len(closestMeterKey)/2+1] == 0 {
			closestMeter, closestMeterName, perfectMatch, tasbeeghOazala = evalMuqatta(sliceCombination, combinationKey, closestMeter, closestMeterName, perfectMatch)
		}
	}

	return dist, closestMeter, closestMeterName, closestMeterKey, perfectMatch, tasbeeghOazala
}

func closestMeter(lineWeights [][][]string, combinedWeights [][]string, bestMatch []string, words [][]string, dict map[string][]string, meters map[string]string, meterNames map[string]string) (string, string, string, []string, []string, bool, int) {

	afaeel := ""
	name := ""
	perfectMatch := false
	sliceCombination := []string{}
	dist := 0
	temp := 0
	index := 0
	meterClosest := ""
	nameClosest := ""
	closestMeterKey := ""
	tasbeeghOazala := 0
	keyClosest := ""
	switch combinedWeights {
	case nil:

		closestMeterKey = ""
		meterClosest = "تمام الفاظ کی شناخت کیے بغیر افاعیل کا تعین ممکن نہیں"
		nameClosest = "تمام الفاظ کی شناخت کیے بغیر بحر کا تعین ممکن نہیں"
		combinedWeights = append(combinedWeights, []string{"تمام الفاظ کی شناخت کیے بغیر تقطیع ممکن نہیں"})
		sliceCombination = []string{""}

	default:

		for i := range combinedWeights {
			combinationKey, sliceCombination := genKey(combinedWeights[i])
			dist, afaeel, name, closestMeterKey, perfectMatch, tasbeeghOazala = meterMatch(combinationKey, bestMatch, sliceCombination, meters, meterNames)
			if i == 0 || dist < temp {
				temp = dist
				index = i
				meterClosest = afaeel
				nameClosest = name
				keyClosest = closestMeterKey
			}
			if perfectMatch == true {
				break
			}
		}
	}

	return meterClosest, nameClosest, keyClosest, sliceCombination, combinedWeights[index], perfectMatch, tasbeeghOazala
}

func checkTasbeeghAzala(combinationKey string) (string, int) {
	tasbeeghOazala := 0
	if strings.HasSuffix(combinationKey, "1") {
		switch {
		case strings.HasSuffix(combinationKey, "0101"):
			tasbeeghOazala = 1
		case strings.HasSuffix(combinationKey, "1101"):
			tasbeeghOazala = 2
		}
		combinationKey = strings.TrimSuffix(combinationKey, "1")
	}
	return combinationKey, tasbeeghOazala
}

func checkMuqatta(closestMeterKey string) bool {
	muqatta := false
	if len(closestMeterKey)%2 == 0 {
		half := len(closestMeterKey) / 2
		halves := chunks(closestMeterKey, half)
		halfQuarter := chunks(halves[0], len(halves[0])/2)
		if halves[0] == halves[1] && halfQuarter[0] != halfQuarter[1] {
			muqatta = true
		}
	}
	return muqatta
}

func evalMuqatta(sliceCombination []string, combinationKey string, closestMeter string, closestMeterName string, perfectMatch bool) (string, string, bool, int) {
	count := 0
	tasbeeghOazala := 0
	for i := range sliceCombination {
		for j := range sliceCombination[i] {
			count++
			if count == len(combinationKey)/2-1 && strings.HasSuffix(sliceCombination[i], "1") && len(sliceCombination[i]) > 1 { //hasSuffox makes sure that the 1 isn't due to a word's first letter
				perfectMatch = true
				_ = j //not using j, only the count matters
				halves := chunks(combinationKey, len(combinationKey)/2+1)
				switch {
				case strings.HasSuffix(halves[0], "0101"):
					closestMeterName += " مسبغ"
					splitAfaeel := strings.Split(closestMeter, " ")
					splitAfaeel[1] = strings.TrimSuffix(splitAfaeel[1], "ن")
					splitAfaeel[1] += "ان"
					closestMeter = strings.Join(splitAfaeel, " ")
					tasbeeghOazala = 3
				case strings.HasSuffix(halves[0], "1101"):
					closestMeterName += " مذال"
					splitAfaeel := strings.Split(closestMeter, " ")
					splitAfaeel[1] = strings.TrimSuffix(splitAfaeel[1], "ن")
					splitAfaeel[1] += "ان"
					closestMeter = strings.Join(splitAfaeel, " ")
					tasbeeghOazala = 4

				}
			}
		}
	}

	return closestMeter, closestMeterName, perfectMatch, tasbeeghOazala
}

func taskeenEAusat(closestMeterKey string, meters map[string]string, meterNames map[string]string) ([]string, []string, []string, bool) {

	ramalMakhboonMusammanFamily := closestMeterKey == "10110101110101110101110" || closestMeterKey == "10110101110101110101010" || closestMeterKey == "1110101110101110101110" || closestMeterKey == "1110101110101110101010"
	ramalMakhboonMurabbaFamily := closestMeterKey == "10110101110101110" || closestMeterKey == "10110101110101010" || closestMeterKey == "1110101110101110" || closestMeterKey == "1110101110101010"
	khafifMusaddasFamily := closestMeterKey == "10110101101101110" || closestMeterKey == "10110101101101010" || closestMeterKey == "1110101101101110" || closestMeterKey == "1110101101101010"
	mujattasMusammanFamily := closestMeterKey == "1101101110101101101110" || closestMeterKey == "1101101110101101101010"
	hazajAkhrabMusaddasFamily := closestMeterKey == "1010111011011010" || closestMeterKey == "1010101011011010"

	taskeen := ramalMakhboonMurabbaFamily || ramalMakhboonMusammanFamily || khafifMusaddasFamily || mujattasMusammanFamily || hazajAkhrabMusaddasFamily
	closestMeterKeys := []string{}
	closestMeterKeys = append(closestMeterKeys, closestMeterKey)
	if ramalMakhboonMusammanFamily {

		switch closestMeterKey {
		case "10110101110101110101110":
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101010")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101010")
		case "10110101110101110101010":
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101010")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101010")
		case "1110101110101110101110":
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101010")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101010")
		case "1110101110101110101010":
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101010")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101010")
			closestMeterKeys = append(closestMeterKeys, "1110101110101110101110")
			closestMeterKeys = append(closestMeterKeys, "10110101110101110101110")
		}
	} else if khafifMusaddasFamily {

		switch closestMeterKey {
		case "10110101101101110":
			{
				closestMeterKeys = append(closestMeterKeys, "10110101101101110")
				closestMeterKeys = append(closestMeterKeys, "10110101101101010")
				closestMeterKeys = append(closestMeterKeys, "1110101101101110")
				closestMeterKeys = append(closestMeterKeys, "1110101101101010")
			}
		case "10110101101101010":
			{
				closestMeterKeys = append(closestMeterKeys, "10110101101101110")
				closestMeterKeys = append(closestMeterKeys, "1110101101101110")
				closestMeterKeys = append(closestMeterKeys, "1110101101101010")
			}
		case "1110101101101110":
			{
				closestMeterKeys = append(closestMeterKeys, "1110101101101010")
				closestMeterKeys = append(closestMeterKeys, "10110101101101110")
				closestMeterKeys = append(closestMeterKeys, "10110101101101010")
			}
		case "1110101101101010":
			{
				closestMeterKeys = append(closestMeterKeys, "1110101101101110")
				closestMeterKeys = append(closestMeterKeys, "10110101101101110")
				closestMeterKeys = append(closestMeterKeys, "10110101101101010")
			}
		}
	} else if ramalMakhboonMurabbaFamily {

		switch closestMeterKey {
		case "10110101110101110":
			{
				closestMeterKeys = append(closestMeterKeys, "10110101110101010")
				closestMeterKeys = append(closestMeterKeys, "1110101110101110")
				closestMeterKeys = append(closestMeterKeys, "1110101110101010")
			}
		case "10110101101101010":
			{
				closestMeterKeys = append(closestMeterKeys, "10110101110101110")
				closestMeterKeys = append(closestMeterKeys, "1110101110101110")
				closestMeterKeys = append(closestMeterKeys, "1110101110101010")
			}
		case "1110101101101110":
			{
				closestMeterKeys = append(closestMeterKeys, "10110101101101010")
				closestMeterKeys = append(closestMeterKeys, "10110101110101110")
				closestMeterKeys = append(closestMeterKeys, "1110101110101010")
			}
		case "1110101101101010":
			{
				closestMeterKeys = append(closestMeterKeys, "1110101101101110")
				closestMeterKeys = append(closestMeterKeys, "10110101110101110")
				closestMeterKeys = append(closestMeterKeys, "10110101101101010")
			}
		}
	} else if mujattasMusammanFamily {

		switch closestMeterKey {
		case "1101101110101101101110":
			{
				closestMeterKeys = append(closestMeterKeys, "1101101110101101101010")
			}
		case "1101101110101101101010":
			{
				closestMeterKeys = append(closestMeterKeys, "1101101110101101101110")
			}

		}
	} else if hazajAkhrabMusaddasFamily {

		switch closestMeterKey {
		case "1010111011011010":
			{
				closestMeterKeys = append(closestMeterKeys, "1010101011011010")
			}
		case "1010101011011010":
			{
				closestMeterKeys = append(closestMeterKeys, "1010111011011010")
			}
		}
	}
	closestMeters, closestMeterNames := []string{}, []string{}
	if ramalMakhboonMurabbaFamily || ramalMakhboonMusammanFamily || khafifMusaddasFamily || hazajAkhrabMusaddasFamily || mujattasMusammanFamily {

		for i := range closestMeterKeys {
			closestMeters = append(closestMeters, meters[closestMeterKeys[i]])
			closestMeterNames = append(closestMeterNames, meterNames[closestMeterKeys[i]])
		}

	} else {
		closestMeters = append(closestMeters, meters[closestMeterKey])
		closestMeterNames = append(closestMeterNames, meterNames[closestMeterKey])
	}
	return closestMeterKeys, closestMeters, closestMeterNames, taskeen
}
