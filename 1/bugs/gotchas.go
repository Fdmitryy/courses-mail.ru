package main

import (
	"sort"
	"strconv"
	"strings"
)

/*
	сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go

	IntSliceToString
	MergeSlices
	GetMapValuesSortedByKey
*/

func IntSliceToString(arr []int) (result string) {
	var str []string
	for _, v := range arr {
		str = append(str, strconv.Itoa(v))
	}
	result = strings.Join(str, "")
	return result
}

func MergeSlices(floatArr []float32, intArr []int32) (result []int) {
	for _, v := range floatArr {
		result = append(result, int(v))
	}
	for _, v := range intArr {
		result = append(result, int(v))
	}
	return result
}

func GetMapValuesSortedByKey(in map[int]string) (result []string) {
	var keys []int
	for k := range in {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		result = append(result, in[k])
	}
	return result
}
