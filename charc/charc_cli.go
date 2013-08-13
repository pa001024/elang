package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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
	count int64

	m  map[rune]int
	ma []rune
}

func NewNormalCounter() *NormalCounter {
	return &NormalCounter{
		m:  make(map[rune]int),
		ma: make([]rune, 0, 0x9fa5-0x4e00),
	}
}
func (this *NormalCounter) ReadAll(in io.Reader) {
	r := bufio.NewReader(in)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			break
		}
		if _, ok := this.m[ru]; ru >= 0x4e00 && ru <= 0x9fa5 {
			if !ok {
				this.m[ru] = 0
				this.ma = append(this.ma, ru)
			}
			this.count++
			this.m[ru]++
		}
	}
}
func (this *NormalCounter) Output(out io.Writer, mutiline bool) {
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
	count int64
	// charsets []Charset
	m  map[rune]int
	rm map[int][]rune
	ra []int
}

func NewSortCounter() *SortCounter { return &SortCounter{m: make(map[rune]int)} }
func (this *SortCounter) ReadAll(in io.Reader) {
	r := bufio.NewReader(in)
	for {
		ru, _, err := r.ReadRune()
		if err != nil {
			break
		}
		if _, ok := this.m[ru]; ru >= 0x4e00 && ru <= 0x9fa5 {
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

func parseflags() (in string, out string, sort, mutiline, count bool) {
	i := flag.String("i", "stdin", "the file you want to use")
	o := flag.String("o", "stdout", "the file you want to output")
	s := flag.Bool("s", false, "sort by use times")
	l := flag.Bool("l", false, "output mutiline text")
	c := flag.Bool("c", false, "print char count")
	h := flag.Bool("h", false, "show help")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	return *i, *o, *s, *l, *c
}

func main() {
	in, out, isSort, isMutiline, isCount := parseflags()
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
		co = NewSortCounter()
	} else {
		co = NewNormalCounter()
	}
	co.ReadAll(fin)
	co.Output(fout, isMutiline)
	if isCount {
		if out == "stdout" || out == "stderr" {
			fmt.Println()
		}
		fmt.Println("[All]", co.AllCount(), "[Chars]", co.Counted())
	}
}
