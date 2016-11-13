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

func (w *ApiWriter) WriteLen(word string) error {
    l := uint32(len(word))
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, l)
    if l > 0x80 {
        buf[2] = buf[2] | 0x80
    }
    var err error
    if l > 0x80 {
        err = w.w.WriteByte(buf[2])
    }
    err = w.w.WriteByte(buf[3])
    return err
}

func (w *ApiWriter) WriteWord(word string) (int, error) {
    w.WriteLen(word)
    i, err := w.w.Write([]byte(word))
    w.w.Flush()
    return i, err
}