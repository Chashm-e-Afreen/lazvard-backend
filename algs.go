package main

import (
	"strings"
)

func cartN(a ...[]string) [][]string { //Cartesian product of n sets
	c := 1
	for _, a := range a {
		c *= len(a)
	}
	if c == 0 {
		return nil
	}
	p := make([][]string, c)
	b := make([]string, c*len(a))
	n := make([]uint64, len(a))
	s := 0
	for i := range p {
		e := s + len(a)
		pi := b[s:e]
		p[i] = pi
		s = e
		for j, n := range n {
			pi[j] = a[j][n]
		}
		for j := len(n) - 1; j >= 0; j-- {
			n[j]++
			if n[j] < uint64(len(a[j])) {
				break
			}
			n[j] = 0
		}
	}
	return p
}

//chunks splits string into n halves
func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

func removeMuqattaWasl(lineWeights [][]string) [][]string {
	index := []int{-1, -1}
	for i := range lineWeights {
		for j := range lineWeights[i] {
			if strings.HasPrefix(lineWeights[i][j], "0") {
				index[0] = i
				index[1] = j
			}
		}
	}
	i := index[0]
	j := index[1]
	lineWeights[i] = append(lineWeights[i][:j-1])
	lineWeights[i-1] = append(lineWeights[i-1][:len(lineWeights[i-1])/2-1])

	return lineWeights
}

func mostFrequent(closestMeterKey []string, perfectMatch []bool) string {

	counterMap := map[string]int{}

	for i := range closestMeterKey {

		if perfectMatch[i] {
			counterMap[closestMeterKey[i]] += 5 //this is a makeshift botch
		}
		counterMap[closestMeterKey[i]]++
		//if the same value gets repeated the count increases at that index
	}

	res := ""
	maxCount := 0
	for i, j := range counterMap {

		if maxCount < j {
			res = i
			maxCount = j
		}
	}

	return res
}

//Contains check if a string exists in a slice of strings
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

//ContainsBool check if a bool exists in a slice of bools
func ContainsBool(a []bool, x bool) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getUrduNumerals(closestScansion []string) []string {
	for i := range closestScansion {
		if !strings.HasPrefix(closestScansion[i], "ان الفاظ") && !strings.HasPrefix(closestScansion[i], "تمام") {
			closestScansion[i] = strings.Replace(closestScansion[i], "1", "۱", -1)
			closestScansion[i] = strings.Replace(closestScansion[i], "0", "۰", -1)
			// closestScansion[i] = reverse(closestScansion[i])
		}
	}

	return closestScansion
}
func copyString2d(src [][]string) [][]string {
	dst := make([][]string, len(src))
	for i := range src {
		dst[i] = copyString1d(src[i])
	}
	return dst
}

func copyString1d(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func copyString3d(src [][][]string) [][][]string {
	dst := make([][][]string, len(src))
	for i := range src {
		dst[i] = copyString2d(src[i])
	}
	return dst
}
