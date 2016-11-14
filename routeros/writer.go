package routeros

import (
	"bufio"
	"encoding/binary"
	"io"
)

type ApiWriter struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *ApiWriter {
	return &ApiWriter{bufio.NewWriter(w)}
}

// Algorithm for length is at http://wiki.mikrotik.com/wiki/Manual:API#Protocol
func (w *ApiWriter) WriteLen(l uint32) error {
	thresholds := []uint32{0, 0x80, 0x4000, 0x200000, 0x10000000}
	masks := []byte{0, 0x80, 0xC0, 0xE0, 0xF0}

	buf := make([]byte, 5)
	binary.BigEndian.PutUint32(buf[1:], l)

	var err error
	masked := false
	for i := 4; i >= 0; i-- {
		if l >= thresholds[i] {
			if !masked {
				buf[4-i] = buf[4-i] | masks[i]
				masked = true
			}
			err = w.w.WriteByte(buf[4-i])
		}
	}
	return err
}

func (w *ApiWriter) WriteWord(word string) (int, error) {
	w.WriteLen(uint32(len(word)))
	i, err := w.w.Write([]byte(word))
	w.w.Flush()
	return i, err
}
