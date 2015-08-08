// go run find_unknown.go|charc -p>unknown_data.txt
package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	bKnownCN, _ := ioutil.ReadFile("known_data.txt")
	knownCN := []rune(string(bKnownCN))
	imap := make(map[rune]int)
	lstr := len(knownCN)
	for i := 0; i < lstr; i++ {
		imap[knownCN[i]] = 1
	}
	data, _ := ioutil.ReadFile("jpconv_lib.txt")
	lib := []rune(string(data))
	llib := len(lib)
	outarr := make([]rune, 0, 1000)
	for i := 0; i < llib; i++ {
		if _, ok := imap[lib[i]]; !ok && lib[i] >= 0x4e00 && lib[i] <= 0x9fa5 {
			outarr = append(outarr, lib[i])
		}
	}
	lout := len(outarr)
	for i := 0; i < lout; i++ {
		fmt.Print(string(outarr[i]))
	}
}
