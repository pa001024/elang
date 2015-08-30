// go run find_unknown.go|charc -p>unknown_data.txt

// new solution (recommand)
/*/
cat known_chs.txt known_cht.txt known_extra.txt > known_data.txt
cat jpconv_lib.txt jpconv_lib_extra.txt | charsub -b known_data.txt | charc -p > unknown_data.txt
rm known_data.txt
//*/
// fix
/*/
mv known_chs.txt known_chs_tmp.txt
charc -p < known_chs_tmp.txt > known_chs.txt
mv known_cht.txt known_cht_tmp.txt
charc -p < known_cht_tmp.txt > known_cht.txt
rm *_tmp.txt

//*/

package main

import (
	"fmt"
	"io/ioutil"
)

func p(inFile, filterFile string) {
	bKnownCN, _ := ioutil.ReadFile(filterFile)
	knownCN := []rune(string(bKnownCN))
	imap := make(map[rune]int)
	lstr := len(knownCN)
	for i := 0; i < lstr; i++ {
		imap[knownCN[i]] = 1
	}
	data, _ := ioutil.ReadFile(inFile)
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

func main() {
	p("jpconv_lib.txt", "known_data.txt")
	p("jpconv_lib_extra.txt", "known_data.txt")
}
