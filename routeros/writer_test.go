package routeros

import (
	"bytes"
	"encoding/binary"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWriter(t *testing.T) {
	Convey("ApiWriter", t, func() {
		buf := new(bytes.Buffer)
		w := NewWriter(buf)
		Convey("for short word", func() {
			w.WriteWord("test")
			Convey("should compute correct word length", func() {
				So(buf.Bytes()[0], ShouldEqual, byte(4))
			})
			Convey("should write correct word", func() {
				So(string(buf.Bytes()[1:]), ShouldEqual, "test")
			})
		})
		Convey("for long word", func() {
			lwbuf := new(bytes.Buffer)
			for i := 'a'; i <= 'z'; i++ {
				for j := 0; j < 10; j++ {
					lwbuf.WriteByte(byte(i))
				}
			}
			w.WriteWord(lwbuf.String())
			Convey("should compute correct word length (260 | 0x8000)", func() {
				So(binary.BigEndian.Uint16(buf.Bytes()[0:2]), ShouldEqual, 260|0x8000)
			})
			Convey("should read correct word", func() {
				So(string(buf.Bytes()[2:]), ShouldEqual, lwbuf.String())
			})
		})
	})
}
