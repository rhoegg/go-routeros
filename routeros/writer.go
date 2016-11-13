package routeros

import (
    "bufio"
    "io"
)

type ApiWriter struct {
    w *bufio.Writer
}

func NewWriter(w io.Writer) *ApiWriter {
    return &ApiWriter{bufio.NewWriter(w)}
}

func (w *ApiWriter) WriteLen(word string) error {
    err := w.w.WriteByte(byte(len(word)))
    return err
}

func (w *ApiWriter) WriteWord(word string) (int, error) {
    w.WriteLen(word)
    i, err := w.w.Write([]byte(word))
    w.w.Flush()
    return i, err
}