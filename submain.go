package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	_ "github.com/lib/pq"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

//DefaultOptionsNew is a set of options suitable for string matching for islah
var DefaultOptionsNew levenshtein.Options = levenshtein.Options{
	InsCost: 5,
	DelCost: 5,
	SubCost: 1,
	Matches: levenshtein.IdenticalRunes,
}

func debug() {
	temp1 := "کسی نے پوچھا کہ حسن کیا ہے"
	temp5 := "ہم نے تیری مثال دے ڈالی"
	newline := "\n"
	// temp2 := "ایک سودا زیاں سے اٹھتا ہے"
	// temp3 := "پئے بیٹھا ہُوں جوش! علم و نظر کے سینکڑوں قُلزم "
	// temp4 := "نکلتے ہیں کبھی تو چاندنی سے دُھوپ کے لشکر"
	// temp4 := "عاشق اگر ہوئے تھے ناز و غرور کیا تھا"

	words, closestScansion, closestMeterKeys, closestMeters, closestMeterNames, problematicWords, ravaniScore := test(temp1 + newline + temp5)

	fmt.Println(words)
	fmt.Println(closestScansion)
	fmt.Println(closestMeterKeys)
	fmt.Println(closestMeters)
	fmt.Println(closestMeterNames)
	fmt.Println(problematicWords)
	fmt.Println(ravaniScore)
}

func test(input string) ([][]string, [][]string, [][]string, []string, []string, [][]bool, []int) {
	// const user string = "postgres"
	// const DSN string = "host=localhost port=5432 user=postgres dbname=myDatabaseLughat sslmode=disable"

	dict := fetchDictFromFile()

	fmt.Print(input)
	temp := removeNuisances(input)
	if len(temp) == 0 {
		return nil, nil, nil, nil, nil, nil, []int{0}
	}
	words := splitInput(temp)

	lineWeights := lineScansion(words, dict)
	fmt.Println(lineWeights)

	combinedWeights := [][][]string{}

	combinedWeights = genAllWeightCombinations(lineWeights)

	afaeel, name, closestMeterKeys, sliceCombination, closestScansion, perfectMatch := make([]string, len(combinedWeights)), make([]string, len(combinedWeights)), make([]string, len(combinedWeights)), make([][]string, len(combinedWeights)), make([][]string, len(combinedWeights)), make([]bool, len(combinedWeights))

	closestMeters := []string{}
	closestMeterNames := []string{}
	problematicWords := [][]bool{}
	tasbeeghOazala := make([]int, len(lineWeights))
	islah := [][]string{}
	ravaniScore := []int{}
	switch { //proceed only if a word is found
	case len(combinedWeights) != 0:
		bestMatch := []string{}

		for i := range combinedWeights {
			afaeel[i], name[i], closestMeterKeys[i], sliceCombination[i], closestScansion[i], perfectMatch[i], tasbeeghOazala[i] = closestMeter(lineWeights, combinedWeights[i], bestMatch, words, dict, meters, meterNames)
		}

		universalKey := mostFrequent(closestMeterKeys, perfectMatch)
		universalKeys, universalAfaeel, universalNames, taskeen := taskeenEAusat(universalKey, meters, meterNames) //check if the behr belongs to a family
		if !taskeen {
			universalKeys = append(universalKeys, universalKey)
			universalAfaeel = append(universalAfaeel, meters[universalKey])
			universalNames = append(universalNames, meterNames[universalKey])
		}

		for i := range combinedWeights {

			afaeel[i], name[i], closestMeterKeys[i], sliceCombination[i], closestScansion[i], perfectMatch[i], tasbeeghOazala[i] = closestMeter(lineWeights, combinedWeights[i], universalKeys, words, dict, meters, meterNames)
		}

		unknownWords := make([][]string, len(words))

		closestMeterKeys, closestMeters, closestMeterNames = tasbeeghOazalaCheck(tasbeeghOazala, closestScansion, closestMeterKeys, closestMeters, closestMeterNames)
		problematicWords, islah = genIslah(closestScansion, closestMeterKeys, perfectMatch, tasbeeghOazala)

		unknownWords = listUnknownWords(lineWeights, words, islah, unknownWords)
		ravaniScore = ravani(words, closestScansion, perfectMatch)

	default:

	}
	return words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore

}

func individualWordVao(vaoBarayeNaam bool, lineWeights [][][]string, vaoIndex []int) ([][][]string, [][][]string) {
	lineWeightsAdditional := make([][][]string, len(lineWeights))
	temp := [][]string{}
	if vaoBarayeNaam {
		for i := 0; i < len(vaoIndex); i += 2 {
			temp = [][]string{}
			temp = append(temp, lineWeights[vaoIndex[i]][:vaoIndex[i+1]]...)
			temp = append(temp, lineWeights[vaoIndex[i]][vaoIndex[i+1]+1:]...)
			lineWeightsAdditional[vaoIndex[i]] = temp
		}
	}
	return lineWeights, lineWeightsAdditional
}

func genAllWeightCombinations(lineWeights [][][]string) [][][]string {
	combinedWeights := make([][][]string, len(lineWeights))
	for i := range lineWeights {
		temp := cartN(lineWeights[i]...)
		if len(temp) > 0 {
			if len(combinedWeights[i]) == 0 {
				combinedWeights[i] = temp
			}
		} else {
			combinedWeights[i] = nil
		}
	}
	return combinedWeights
}

func tasbeeghOazalaCheck(tasbeeghOazala []int, closestScansion [][]string, closestMeterKeys []string, closestMeters []string, closestMeterNames []string) ([]string, []string, []string) {
	for i := range closestMeterKeys {
		if tasbeeghOazala[i] == 0 {
			closestMeters = append(closestMeters, meters[closestMeterKeys[i]])
			closestMeterNames = append(closestMeterNames, meterNames[closestMeterKeys[i]])
		} else {
			// /closestMeterKeys[i] = closestMeterKeys[i][:len(closestMeterKeys[i])/2] + closestMeterKeys[i][:len(closestMeterKeys[i])/2]
			closestMeters = append(closestMeters, meters[closestMeterKeys[i]])
			closestMeterNames = append(closestMeterNames, meterNames[closestMeterKeys[i]])

		}
		if tasbeeghOazala[i] == 1 {
			closestMeterKeys[i] += "1"
			closestMeterNames[i] += " مسبغ"
			closestMeters[i] = strings.TrimSuffix(closestMeters[i], "ن") + "ان"

		}
		if tasbeeghOazala[i] == 2 {
			closestMeterKeys[i] += "1"
			closestMeterNames[i] += " مذال"
			closestMeters[i] = strings.TrimSuffix(closestMeters[i], "ن") + "ان"

		}
		if tasbeeghOazala[i] == 3 {

			closestMeterKeys[i] = closestMeterKeys[i][:len(closestMeterKeys[i])/2]
			closestMeterKeys[i] = closestMeterKeys[i] + "1" + closestMeterKeys[i]

			closestMeterNames[i] += " مسبغ"
			splitAfaeel := strings.Split(closestMeters[i], " ")
			splitAfaeel[1] = strings.TrimSuffix(splitAfaeel[1], "ن")
			splitAfaeel[1] += "ان"
			closestMeters[i] = strings.Join(splitAfaeel, " ")
			if strings.HasSuffix(closestScansion[i][len(closestScansion[i])-1], "1") {
				closestMeterKeys[i] += "1"
				closestMeters[i] = strings.TrimSuffix(closestMeters[i], "ن") + "ان"
			}
		}
		if tasbeeghOazala[i] == 4 {

			closestMeterKeys[i] = closestMeterKeys[i][:len(closestMeterKeys[i])/2]
			closestMeterKeys[i] = closestMeterKeys[i] + "1" + closestMeterKeys[i]
			closestMeterNames[i] += " مذال"
			splitAfaeel := strings.Split(closestMeters[i], " ")
			splitAfaeel[1] = strings.TrimSuffix(splitAfaeel[1], "ن")
			splitAfaeel[1] += "ان"
			closestMeters[i] = strings.Join(splitAfaeel, " ")
			if strings.HasSuffix(closestScansion[i][len(closestScansion[i])-1], "1") {
				closestMeterKeys[i] += "1"
				closestMeters[i] = strings.TrimSuffix(closestMeters[i], "ن") + "ان"
			}
		}
	}
	return closestMeterKeys, closestMeters, closestMeterNames
}

func listUnknownWords(lineWeights [][][]string, words [][]string, closestMeterKeys [][]string, unknownWords [][]string) [][]string {

	for i := range lineWeights {
		for j := range lineWeights[i] {
			if len(lineWeights[i][j]) == 0 {
				unknownWords[i] = append(unknownWords[i], words[i][j])
			}
		}
	}

	for i := range unknownWords {

		notFound := "ان الفاظ کی شناخت نہ کی جا سکی" + ": "
		for j := range unknownWords[i] {
			notFound += " " + unknownWords[i][j]
		}
		if !strings.HasPrefix(closestMeterKeys[i][0], "1") && !strings.HasPrefix(closestMeterKeys[i][0], "0") {
			closestMeterKeys[i][0] = notFound
		}
	}

	return unknownWords
}

func ravani(words [][]string, closestScansion [][]string, perfectMatch []bool) []int {
	ravaniScore := make([]int, len(words))

	for i := range words {
		ravaniScore[i] = 10
		switch {
		case !strings.HasPrefix(closestScansion[i][0], "تمام") && perfectMatch[i] == true:
			for j := range words[i] {
				//
				if j > 0 && strings.HasSuffix(closestScansion[i][j-1], "1") {
					lastLetter, _ := utf8.DecodeLastRuneInString(words[i][j-1])
					firstLetter, _ := utf8.DecodeRuneInString(words[i][j])
					if firstLetter == lastLetter {
						ravaniScore[i]--
					}
				}

				if strings.HasSuffix(words[i][j], "ِ") {
					if j > 0 && strings.HasSuffix(closestScansion[i][j-1], "1") {

						lastLetter, _ := utf8.DecodeLastRuneInString(words[i][j-1][:len(words[i][j-1])-2])
						firstLetter, _ := utf8.DecodeRuneInString(words[i][j])
						if firstLetter == lastLetter {
							ravaniScore[i]--
						}
					}
				}
				hindiUlAsl := strings.Contains(words[i][j], "ٹ") || strings.Contains(words[i][j], "ڈ") || strings.Contains(words[i][j], "ڑ") || strings.Contains(words[i][j], "ھ")

				if j < len(words[i])-1 && (strings.HasSuffix(words[i][j], "ا") || strings.HasSuffix(words[i][j], "ی") || strings.HasSuffix(words[i][j], "ے") || strings.HasSuffix(words[i][j], "ں")) {
					if j > 0 && strings.HasSuffix(closestScansion[i][j-1], "1") {

						lastLetter, _ := utf8.DecodeLastRuneInString(words[i][j-1][:len(words[i][j-1])-2])
						firstLetter, _ := utf8.DecodeRuneInString(words[i][j])
						if firstLetter == lastLetter {
							ravaniScore[i]--
						}

						switch {
						case strings.HasSuffix(closestScansion[i][j], "1") && (hindiUlAsl || len(words[i][j]) > 5 || strings.HasSuffix(words[i][j], "ے")):
							ravaniScore[i]--
							if strings.HasSuffix(words[i][j-1], "اں") || strings.HasSuffix(words[i][j-1], "یں") || strings.HasSuffix(words[i][j-1], "وں") && len(words[i][j-1]) > 4 {
								ravaniScore[i]--
							}
							if strings.HasSuffix(words[i][j-1], "ا") && len(words[i][j-1]) > 4 {
								ravaniScore[i]--
							}
						default:
							if j > 0 && strings.HasSuffix(closestScansion[i][j], "1") {
								if strings.HasSuffix(words[i][j], "ے") {
									ravaniScore[i]--
								} else {
									ravaniScore[i] -= 2
								}
							}
						}

					}
				}
			}

		default:
			ravaniScore[i] = 0
		}
	}

	return ravaniScore

}

func genIslah(scansion [][]string, keys []string, perfectMatch []bool, tasbeeghOazala []int) ([][]bool, [][]string) {
	closestScansion := copyString2d(scansion)
	problematicWords := make([][]bool, len(closestScansion))
	closestMeterKeys := copyString1d(keys)
	islah := closestScansion

	for i := range closestScansion {
		if !strings.HasPrefix(scansion[i][0], "تمام الفاظ") {
			if !perfectMatch[i] {
				problematicWords[i] = make([]bool, len(closestScansion[i]))
				if tasbeeghOazala[i] != 0 { //for muqatta bahoor
					closestMeterKeys[i] = closestMeterKeys[i][:len(closestMeterKeys[i])/2]
					closestMeterKeys[i] = closestMeterKeys[i] + "1" + closestMeterKeys[i]
				}
				temp := strings.Join(closestScansion[i][:], "")
				editScript := levenshtein.EditScriptForStrings([]rune(temp), []rune(closestMeterKeys[i]), DefaultOptionsNew) //levenshtein script to optimise problematic words
				count := 0
				for j := range islah[i] {
					operations := 0

					for k := count + operations; k < count+len(islah[i][j])+operations; k++ {
						switch editScript[k] {
						case 0:
							operations++
							problematicWords[i][j] = true
						case 1:
							operations--
							problematicWords[i][j] = true
						case 2:
							problematicWords[i][j] = true
						}
					}

					islah[i][j] = closestMeterKeys[i][count : count+len(islah[i][j])+operations]

					count += len(islah[i][j])
				}
				for t := range islah[i] {
					if t > 0 && strings.HasPrefix(islah[i][t], "0") {
						if len(islah[i][t-1]) > 2 { //avoids a situation when the suggested weight for a word is 1
							islah[i][t] = "1" + islah[i][t]
							islah[i][t-1] = strings.TrimSuffix(islah[i][t-1], "1")
							problematicWords[i][t-1] = true
							if islah[i][t] == scansion[i][t] {
								problematicWords[i][t] = false
							}
						} else {
							islah[i][t] = strings.TrimPrefix(islah[i][t], "0")
							islah[i][t-1] += "0"
							problematicWords[i][t-1] = true
							if islah[i][t] == scansion[i][t] {
								problematicWords[i][t] = false
							}
						}
					}

				}

			} else {
				for j := range closestScansion[i] {
					problematicWords[i] = make([]bool, len(closestScansion[i]))
					problematicWords[i][j] = false
				}
			}
		} else {
			problematicWords[i] = []bool{true}
		}
	}

	return problematicWords, islah
}
