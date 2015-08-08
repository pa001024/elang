package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
)

type Charset struct {
	From, To rune
}

type Counter interface {
	ReadAll(in io.Reader)                // 读取全部
	Output(out io.Writer, mutiline bool) // 输出
	AllCount() int64                     // 全部字符数量
	Counted() int                        // 计算进的数量
}

// 不排序字符统计器
type NormalCounter struct {
	count       int64
	isAllChar   bool
	isSortASCII bool
	m           map[rune]int
	ma          []rune
}

func NewNormalCounter(allChar, sortASCII bool) *NormalCounter {
	maxcap := 0x9fa5 - 0x4e00
	if allChar {
		maxcap = 0xffff
	}
	return &NormalCounter{
		m:           make(map[rune]int),
		ma:          make([]rune, 0, maxcap),
		isAllChar:   allChar,
		isSortASCII: sortASCII,
	}
}
func (this *NormalCounter) ReadAll(in io.Reader) {
	r := bufio.NewReader(in)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			break
		}
		if _, ok := this.m[ru]; this.isAllChar || ru >= 0x4e00 && ru <= 0x9fa5 {
			if !ok {
				this.m[ru] = 0
				this.ma = append(this.ma, ru)
			}
			this.count++
			this.m[ru]++
		}
	}
}
func (this *NormalCounter) SortASCII() {
	sorter := &runeSorter{this.ma}
	sort.Sort(sorter)
}

type runeSorter struct{ runes []rune }

func (s *runeSorter) Len() int           { return len(s.runes) }
func (s *runeSorter) Swap(i, j int)      { s.runes[i], s.runes[j] = s.runes[j], s.runes[i] }
func (s *runeSorter) Less(i, j int) bool { return s.runes[i] < s.runes[j] }

func (this *NormalCounter) Output(out io.Writer, mutiline bool) {
	if this.isSortASCII {
		this.SortASCII()
	}
	w := bufio.NewWriter(out)
	for _, v := range this.ma {
		v2 := this.m[v]
		if mutiline {
			fmt.Fprintf(w, "%s : %v\n", string(v), v2)
		} else {
			w.WriteRune(v)
		}
	}
	w.Flush()
}
func (this *NormalCounter) AllCount() int64 { return this.count }
func (this *NormalCounter) Counted() int    { return len(this.ma) }

// 排序字符统计器
type SortCounter struct {
	count     int64
	isAllChar bool
	m         map[rune]int
	rm        map[int][]rune
	ra        []int
}

func NewSortCounter(allChar bool) *SortCounter {
	return &SortCounter{m: make(map[rune]int), isAllChar: allChar}
}
func (this *SortCounter) ReadAll(in io.Reader) {
	r := bufio.NewReader(in)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			break
		}
		if _, ok := this.m[ru]; this.isAllChar || ru >= 0x4e00 && ru <= 0x9fa5 {
			if !ok {
				this.m[ru] = 0
			}
			this.count++
			this.m[ru]++
		}
	}
}
func (this *SortCounter) Sort() {
	this.rm = make(map[int][]rune)
	for i, v := range this.m {
		if _, ok := this.rm[v]; !ok {
			this.rm[v] = make([]rune, 0, 2)
		}
		this.rm[v] = append(this.rm[v], i)
	}
	this.ra = make([]int, 0, len(this.rm))
	for i, _ := range this.rm {
		this.ra = append(this.ra, i)
	}
	sort.Ints(this.ra)
}
func (this *SortCounter) Output(out io.Writer, mutiline bool) {
	this.Sort()
	w := bufio.NewWriter(out)
	for i := len(this.ra) - 1; i >= 0; i-- {
		v := this.ra[i]
		for _, v2 := range this.rm[v] {
			if mutiline {
				fmt.Fprintf(w, "%s : %v\n", string(v2), v)
			} else {
				w.WriteRune(v2)
			}
		}
	}
	w.Flush()
}
func (this *SortCounter) AllCount() int64 { return this.count }
func (this *SortCounter) Counted() int    { return len(this.m) }

func parseflags() (in string, out string, sort, sortascii, mutiline, allchar, count, random bool) {
	i := flag.String("i", "stdin", "the file you want to use")
	o := flag.String("o", "stdout", "the file you want to output")
	s := flag.Bool("s", false, "sort by use times")
	p := flag.Bool("p", false, "sort by ASCII")
	l := flag.Bool("l", false, "output mutiline text")
	a := flag.Bool("a", false, "output all char")
	c := flag.Bool("c", false, "print char count")
	h := flag.Bool("h", false, "show help")
	r := flag.Bool("r", false, "random output")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	return *i, *o, *s, *p, *l, *a, *c, *r
}

func main() {
	in, out, isSort, isSortASCII, isMutiline, isAllChar, isCount, isRandom := parseflags()
	var fin io.Reader
	if in == "stdin" {
		fin = os.Stdin
	} else {
		f, err := os.Open(in)
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}
		defer f.Close()
		fin = f
	}
	var fout io.Writer
	if out == "stdout" {
		fout = os.Stdout
	} else if out == "stderr" {
		fout = os.Stderr
	} else {
		fw, err := os.Create(out)
		if err != nil {
			fmt.Println("[ERROR]", err)
		}
		defer fw.Close()
		fout = fw
	}
	var co Counter
	if isSort {
		co = NewSortCounter(isAllChar)
	} else {
		co = NewNormalCounter(isAllChar, isSortASCII)
	}

	// stdin 输入下防止坑爹
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if isCount {
				fmt.Println("\n[All]", co.AllCount(), "[Chars]", co.Counted())
			}
			os.Exit(0)
		}
	}()
	co.ReadAll(fin)
	if isRandom {
		buf := &bytes.Buffer{}
		co.Output(buf, false)
		io.Copy(fout, buf)
	} else {
		co.Output(fout, isMutiline)
	}
	if isCount {
		fmt.Println("\n[All]", co.AllCount(), "[Chars]", co.Counted())
	}
}
