package main

import (
	// "github.com/axgle/mahonia"
	"io"
	"os"

	"bytes"
	"strings"
)

func main() {
	s := strings.Join(os.Args[1:], " ")
	// g2u := mahonia.NewDecoder("gbk")
	// r := g2u.NewReader(bytes.NewBufferString(s))
	io.Copy(os.Stdout, bytes.NewBufferString(s))
}
