package routeros

import (
    . "github.com/smartystreets/goconvey/convey"
    "testing"
    "bytes"
//    "encoding/binary"
)

func TestWriter(t *testing.T) {
    Convey("ApiWriter", t, func() {
        buf := new(bytes.Buffer)
        Convey("for short word", func() {
            w := NewWriter(buf)
            w.WriteWord("test")
            Convey("should compute correct word length", func() {
                So(buf.Bytes()[0], ShouldEqual, byte(4))
            })
            Convey("should read correct word", func() {
                So(string(buf.Bytes()[1:]), ShouldEqual, "test")
            })
        })
    })
}
