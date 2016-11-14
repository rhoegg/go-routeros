package routeros

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		Convey("WriteLen", func() {
			Convey("for single byte length", func() {
				cases := []int{1, 10, 99, 127}
				for _, i := range cases {
					Convey(fmt.Sprintf("%d", i), func() {
						So(writeLenAndReadLength(w, buf, i, 1), ShouldEqual, i)
					})
				}
			})
			Convey("for two byte length", func() {
				cases := []int{0x80, 1000, 0x3FFF}
				for _, i := range cases {
					Convey(fmt.Sprintf("%d (0x%X)", i, i), func() {
						So(writeLenAndReadLength(w, buf, i, 2), ShouldEqual, i|0x8000)
					})
				}
			})
			Convey("for three byte length", func() {
				cases := []int{0x4000, 100000, 0x1FFFFF}
				for _, i := range cases {
					Convey(fmt.Sprintf("%d (0x%X)", i, i), func() {
						So(writeLenAndReadLength(w, buf, i, 3), ShouldEqual, i|0xC00000)
					})
				}
			})
			Convey("for four byte length", func() {
				cases := []int{0x200000, 123123123, 0xFFFFFFF}
				for _, i := range cases {
					Convey(fmt.Sprintf("%d (0x%X)", i, i), func() {
						So(writeLenAndReadLength(w, buf, i, 4), ShouldEqual, i|0xE0000000)
					})
				}
			})
			Convey("for five byte length", func() {
				cases := []int{0x10000000, 888888888, 0xFFFFFFFF}
				for _, i := range cases {
					Convey(fmt.Sprintf("%d (0x%X)", i, i), func() {
						So(writeLenAndReadLength(w, buf, i, 5), ShouldEqual, i)
						So(buf.Bytes()[0], ShouldEqual, 0xF0)
					})
				}
			})
		})
	})
}

func writeLenAndReadLength(w *ApiWriter, buf *bytes.Buffer, v int, bytes int) uint32 {
	w.WriteLen(uint32(v))
	w.w.Flush()
	So(buf.Len(), ShouldEqual, bytes)
	var r uint32
	if bytes <= 4 {
		lenbuf := make([]byte, 4)
		for i := 0; i < bytes; i++ {
			lenbuf[4-bytes+i] = buf.Bytes()[i]
		}
		r = binary.BigEndian.Uint32(lenbuf)
	} else {
		r = binary.BigEndian.Uint32(buf.Bytes()[1:5])
	}
	return r
}
