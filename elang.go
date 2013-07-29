/*
鹅语
===========

  《赤装宁》
   by 弄姥尸
装宁徴宝岚繁勧，
岷崛畑辅逊安刔。
送佩暖锻弌赔仟，
畅酎抹亟赤装宁。

*/

package elang

import (
	"bytes"
	"io"

	"code.google.com/p/mahonia"
)

var (
	// 正转 u | u2j | g2u = u | u2e
	u2g = mahonia.NewEncoder("gbk")
	j2u = mahonia.NewDecoder("euc-jp")

	// 逆转 u | u2g | j2u = u | e2u
	u2j = mahonia.NewEncoder("euc-jp")
	g2u = mahonia.NewDecoder("gbk")
)

type Decoder struct {
	r    io.Reader
	u2jw io.Writer
	g2ur io.Reader
}

func (this *Decoder) Read(p []byte) (n int, err error) {
	io.Copy(this.u2jw, this.r)
	n, err = this.g2ur.Read(p)
	return
}
func NewDecoder(r io.Reader) io.Reader {
	buf := &bytes.Buffer{}
	return &Decoder{r, u2j.NewWriter(buf), g2u.NewReader(buf)}
}

type Encoder struct {
	w    io.Writer
	u2gw io.Writer
	j2ur io.Reader
}

func NewEncoder(w io.Writer) io.WriteCloser {
	buf := &bytes.Buffer{}
	return &Encoder{w, u2g.NewWriter(buf), j2u.NewReader(buf)}
}
func (this *Encoder) Write(p []byte) (n int, err error) {
	n, err = this.u2gw.Write(p)
	io.Copy(this.w, this.j2ur)
	return
}
func (this *Encoder) Close() error {
	return nil
}

// 转换中文到鹅语
func EncodeString(str string) string {
	buf := &bytes.Buffer{}
	e := NewEncoder(buf)
	io.Copy(e, bytes.NewBufferString(str))
	return buf.String()
}

// 转换鹅语到中文
func DecodeString(str string) string {
	buf := &bytes.Buffer{}
	e := NewDecoder(bytes.NewBufferString(str))
	io.Copy(buf, e)
	return buf.String()
}
