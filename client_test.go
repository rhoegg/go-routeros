package routeros

import (
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "bytes"
)

func TestReader(t *testing.T) {
    Convey("ApiReader", t, func() {
        buf := new(bytes.Buffer)
        Convey("for short word", func() {
            buf.WriteByte(byte(4))
            buf.Write([]byte("test"))
            buf.Write([]byte("trash"))
            r := NewReader(buf)
            Convey("should compute correct word length", func() {
                l, _ := r.ReadLen()
                So(l, ShouldEqual, 4)
            })
            Convey("should read correct word", func() {
                w, _ := r.ReadWord()
                So(w, ShouldEqual, "test")
            })
        })
    })
}