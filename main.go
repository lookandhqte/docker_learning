package main

import (
	"fmt"
	"strings"
)

type StringToInt interface {
	StringToIntConv([]string) []int
}

type StringConversation struct {
}

func (sс StringConversation) StringToIntConv(arr []string) []int {
	result := make([]int, 0, len(arr))

	for _, str := range arr {
		result = append(result, strings.Count(str, "x"))
		// result[id] = strings.Count(str, "x")
	}
	return result
}

func main() {

	arr := [...]string{"xxx", "meow", "meow", "i love x", "deepseek"}
	converter := StringConversation{}
	result := converter.StringToIntConv(arr[:])
	fmt.Println(result)
}
