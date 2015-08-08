package main

import (
	"flag"
	"io"
	"os"

	"github.com/pa001024/elang/jpconv"
)

func main() {
	isDecode := flag.Bool("d", false, "Decode?")
	flag.Parse()
	if *isDecode {
		io.Copy(os.Stdout, jpconv.NewDecoder(os.Stdin))
	} else {
		io.Copy(jpconv.NewEncoder(os.Stdout), os.Stdin)
	}
}
