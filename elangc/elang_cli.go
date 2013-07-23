package main

import (
	"flag"
	"github.com/pa001024/elang"
	"io"
	"os"
)

func main() {
	isDecode := flag.Bool("d", false, "Decode?")
	flag.Parse()
	if *isDecode {
		io.Copy(os.Stdout, elang.NewDecoder(os.Stdin))
	} else {
		io.Copy(elang.NewEncoder(os.Stdout), os.Stdin)
	}
}
