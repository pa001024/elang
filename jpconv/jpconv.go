/*
基础简和转换(不转换词组)

以后可能使用大数组优化 不过考虑到目前字符数量相对还是较少 所以还是用map实现
*/

package jpconv

import (
	"bytes"
	"io"
	"unicode/utf8"
)

type Decoder struct {
	r     io.Reader
	inbuf *bytes.Buffer
}

func (this *Decoder) Read(p []byte) (n int, err error) {
	io.Copy(this.inbuf, this.r)
	n, err = this.inbuf.Read(p)

	d := []rune(string(p[:n]))
	for i, v := range d {
		if e, ok := s2j[v]; ok && len(e) > 0 {
			d[i] = e[0]
		}
	}
	out := []byte(string(d))
	if l := len(out); l <= len(p) {
		copy(p, out)
		n = l
	} else {
		copy(p, out[:l])
		n = len(p)
	}
	return
}
func NewDecoder(r io.Reader) io.Reader {
	return &Decoder{r, &bytes.Buffer{}}
}

type Encoder struct {
	w      io.Writer
	inbuf  []byte
	outbuf []byte
}

func NewEncoder(w io.Writer) io.WriteCloser {
	return &Encoder{w, []byte{}, []byte{}}
}
func (this *Encoder) Write(p []byte) (n int, err error) {
	n = len(p)

	if len(this.inbuf) > 0 {
		this.inbuf = append(this.inbuf, p...)
		p = this.inbuf
	}
	if len(this.outbuf) < len(p) {
		this.outbuf = make([]byte, len(p)+10)
	}
	outpos := 0
	for len(p) > 0 {
		rune, size := utf8.DecodeRune(p)
		if rune == 0xfffd && !utf8.FullRune(p) {
			break
		}
		p = p[size:]

		if e, ok := j2s[rune]; ok && len(e) > 0 {
			rune = e[0]
		}
	retry:
		size = utf8.EncodeRune(this.outbuf[outpos:], rune)
		if size == 0 { // 死循环?
			newDest := make([]byte, len(this.outbuf)*2)
			copy(newDest, this.outbuf)
			this.outbuf = newDest
			goto retry
		}
		outpos += size
	}

	this.inbuf = this.inbuf[:0]
	if len(p) > 0 {
		this.inbuf = append(this.inbuf, p...)
	}

	n1, err := this.w.Write(this.outbuf[0:outpos])

	if err != nil && n1 < n {
		n = n1
	}
	return
}
func (this *Encoder) Close() error {
	return nil
}

// 转换和体到简体
func EncodeString(str string) string {
	buf := &bytes.Buffer{}
	e := NewEncoder(buf)
	io.Copy(e, bytes.NewBufferString(str))
	return buf.String()
}

// 转换简体到和体
func DecodeString(str string) string {
	buf := &bytes.Buffer{}
	e := NewDecoder(bytes.NewBufferString(str))
	io.Copy(buf, e)
	return buf.String()
}
