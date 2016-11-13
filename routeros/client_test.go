package routeros

import (
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "bytes"
    "encoding/binary"
)

func TestReader(t *testing.T) {
    Convey("ApiReader", t, func() {
        buf := new(bytes.Buffer)
        Convey("for short word", func() {
            buf.WriteByte(byte(4))
            buf.Write([]byte("test and some trash"))
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

        Convey("for long word", func() {
            lwbuf := new(bytes.Buffer)
            for i := 'a'; i <= 'z'; i++ {
                for j := 0; j < 10; j++ {
                    lwbuf.WriteByte(byte(i))
                }
            }
            bs := make([]byte, 2)
            binary.BigEndian.PutUint16(bs, 0x8000 | 260)
            buf.Write(bs)
            buf.Write(lwbuf.Bytes())
            r := NewReader(buf)
            Convey("should compute correct word length", func() {
                l, _ := r.ReadLen()
                So(l, ShouldEqual, 260)
            })
            Convey("should read correct word", func() {
                w, _ := r.ReadWord()
                So(w, ShouldEqual, lwbuf.String())
            })
        })
    })
}