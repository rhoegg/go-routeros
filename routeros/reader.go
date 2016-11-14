package routeros

import (
	"bufio"
	"io"
	"io/ioutil"
)

type ApiReader struct {
	rd *bufio.Reader
}

func NewReader(rd io.Reader) *ApiReader {
	return &ApiReader{bufio.NewReader(rd)}
}

func (r *ApiReader) ReadWord() (string, error) {
	l, err := r.ReadLen()
	data, err := ioutil.ReadAll(io.LimitReader(r.rd, int64(l)))
	return string(data), err
}

func (r *ApiReader) ReadSentence() ([]string, error) {
	s := []string{}
	w, err := r.ReadWord()
	for len(w) > 0 && err == nil {
		s = append(s, w)
		w, err = r.ReadWord()
	}
	return s, err
}

// Algorithm for length is at http://wiki.mikrotik.com/wiki/Manual:API#Protocol
func (r *ApiReader) ReadLen() (uint32, error) {
	threshold := [...]byte{0, 0x80, 0xC0, 0xE0, 0xF0}

	first, err := r.rd.ReadByte()
	i := uint32(first)
	for t := 0; t < len(threshold); t++ {
		if first < threshold[t+1] {
			// 0x00, 0x8000, 0xC00000, 0xE0000000
			mask := uint32(threshold[t]) << uint32(t*8)
			return i &^ mask, err
		}
		b, _ := r.rd.ReadByte()
		i = i<<8 + uint32(b)
	}
	return i, ErrUnsupportedWordLength
}
