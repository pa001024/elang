package main

import (
	"flag"
	"github.com/pa001024/elang"
	"io"
	"os"

	"bytes"
	"fmt"
)

func main() {
	isDecode := flag.Bool("d", false, "Decode?")
	isView := flag.Bool("v", false, "View?")
	flag.Parse()
	if *isDecode {
		if *isView {
			buf := &bytes.Buffer{}
			io.Copy(buf, elang.NewDecoder(os.Stdin))
			fmt.Printf("%#v\n", buf.Bytes())
			io.Copy(os.Stdout, buf)
		} else {
			io.Copy(os.Stdout, elang.NewDecoder(os.Stdin))
		}
	} else {
		if *isView {
			buf := &bytes.Buffer{}
			io.Copy(elang.NewEncoder(buf), os.Stdin)
			fmt.Printf("%#v\n", buf.Bytes())
			io.Copy(os.Stdout, buf)
		} else {
			io.Copy(elang.NewEncoder(os.Stdout), os.Stdin)
		}
	}
}
