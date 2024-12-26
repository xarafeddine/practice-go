package main

import (
	"fmt"

	"time"
)

func show() {
	for i := 0; i < 2; i++ {
		time.Sleep(1 * 1e9)
		fmt.Println("In the show")
	}
}

func matchingStrings(strings []string, queries []string) []int32 {
	// Write your code here
	hashMap := map[string]int32{}
	for _, str := range queries {
		hashMap[str] = 0
	}

	for _, str := range strings {
		hashMap[str]++
	}
	var result []int32
	for _, v := range queries {
		result = append(result, hashMap[v])
	}

	return result
}
func main() {

	logParser("text.txt", "normalFile", "errFile")

}
