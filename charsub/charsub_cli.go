package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
)

type Filter interface {
	io.Reader
}

// 一般过滤器
type NormalFilter struct {
	r         io.Reader
	inbuf     *bytes.Buffer
	blacklist map[rune]int8
	whitelist map[rune]int8
}

func fillList(filename string, list *map[rune]int8) {
	var rd io.Reader
	if filename == "stdin" {
		rd = os.Stdin
	} else {
		rdf, err := os.Open(filename)
		if err != nil {
			return
		}
		rd = rdf
		defer rdf.Close()
	}
	*list = make(map[rune]int8)
	brd := bufio.NewReader(rd)
	for {
		re, _, err := brd.ReadRune()
		if err != nil {
			break
		}
		(*list)[re] = 1
	}
}

func NewNormalFilter(fin io.Reader, blacklistFile, whitelistFile string) *NormalFilter {
	var blacklist, whitelist map[rune]int8
	if blacklistFile != "" {
		fillList(blacklistFile, &blacklist)
	}
	if whitelistFile != "" {
		fillList(whitelistFile, &whitelist)
	}

	return &NormalFilter{
		r:         fin,
		inbuf:     &bytes.Buffer{},
		blacklist: blacklist,
		whitelist: whitelist,
	}
}
func (this *NormalFilter) Read(p []byte) (n int, err error) {
	io.Copy(this.inbuf, this.r)
	n, err = this.inbuf.Read(p)
	d := []rune(string(p[:n]))
	outd := make([]rune, 0, len(d))
	for _, v := range d {
		inBlack := false
		if this.blacklist != nil {
			_, inBlack = this.blacklist[v]
		}
		inWhite := true
		if this.whitelist != nil {
			_, inWhite = this.whitelist[v]
		}
		if !inBlack && inWhite {
			outd = append(outd, v)
		}
	}
	out := []byte(string(outd))
	if l := len(out); l <= len(p) {
		copy(p, out)
		n = l
	} else {
		copy(p, out[:l])
		n = len(p)
	}
	return
}

func parseflags() (in, out, blacklist, whitelist string) {
	i := flag.String("i", "stdin", "the file you want to use")
	o := flag.String("o", "stdout", "the file you want to output")
	b := flag.String("b", "", "the file(or stdin) you use as blacklist")
	w := flag.String("w", "", "the file(or stdin) you use as whitelist")
	h := flag.Bool("h", false, "show help")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	return *i, *o, *b, *w
}

func main() {
	in, out, blacklist, whitelist := parseflags()
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
	var mf Filter
	mf = NewNormalFilter(fin, blacklist, whitelist)

	// stdin 输入下防止坑爹
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			os.Exit(0)
		}
	}()
	io.Copy(fout, mf)
}
