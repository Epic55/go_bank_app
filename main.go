package main

import (
	"fmt"
	"unicode"
)

func main() {
	var l1 int
	var s1 string = "abcd"
	l1 = len(s1)
	if l1 < 4 || l1 > 25 {
		fmt.Println("F")
	} else {
		fmt.Println("T")
	}
	fmt.Println(string(s1[0]))
	fmt.Println(s1[0:1])
	if unicode.IsLetter(string(s1[0])) {
		fmt.Println("T")
	}
}
