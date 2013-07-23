package main

import (
	"flag"
	"io"
	"os"

	"github.com/pa001024/elang/zhconv"
)

func main() {
	isDecode := flag.Bool("d", false, "Decode?")
	flag.Parse()
	if *isDecode {
		io.Copy(os.Stdout, zhconv.NewDecoder(os.Stdin))
	} else {
		io.Copy(zhconv.NewEncoder(os.Stdout), os.Stdin)
	}
}
