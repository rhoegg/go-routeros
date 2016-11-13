package routeros

import (
	"bytes"
	"encoding/binary"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
			binary.BigEndian.PutUint16(bs, 0x8000|260)
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

		Convey("ReadLen", func() {
			bs := make([]byte, 5)
			for i := 0; i < 5; i++ {
				bs[i] = byte(245) // fill it with trash to start
			}
			Convey("for single byte length", func() {
				Convey("1", func() {
					bs[0] = byte(1)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 1)
				})
				Convey("9", func() {
					bs[0] = byte(9)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 9)
				})
				Convey("100", func() {
					bs[0] = byte(100)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 100)
				})
				Convey("127", func() {
					bs[0] = byte(127)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 127)
				})
			})
			Convey("for two byte length", func() {
				Convey("128 (0x0080)", func() {
					binary.BigEndian.PutUint16(bs[0:2], 0x8000|128)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 128)
				})
				Convey("999 (0x03E7)", func() {
					binary.BigEndian.PutUint16(bs[0:2], 0x8000|999)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 999)
				})
				Convey("16383 (0x3FFF)", func() {
					binary.BigEndian.PutUint16(bs[0:2], 0x8000|16383)
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 16383)
				})
			})
			Convey("for three byte length", func() {
				b32 := make([]byte, 4)
				Convey("16384 (0x4000)", func() {
					binary.BigEndian.PutUint32(b32, 16384)
					copy(bs[0:3], b32[1:4]) // last 3 bytes
					bs[0] |= 0xC0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 16384)
				})
				Convey("1000001 (0x0F4241)", func() {
					binary.BigEndian.PutUint32(b32, 1000001)
					copy(bs[0:3], b32[1:4]) // last 3 bytes
					bs[0] |= 0xC0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 1000001)
				})
				Convey("2097151 (0x1FFFFF)", func() {
					binary.BigEndian.PutUint32(b32, 2097151)
					copy(bs[0:3], b32[1:4]) // last 3 bytes
					bs[0] |= 0xC0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 2097151)
				})
			})
			Convey("for four byte length", func() {
				b32 := make([]byte, 4)
				Convey("2097152 (0x200000)", func() {
					binary.BigEndian.PutUint32(b32, 2097152)
					copy(bs, b32[0:4])
					bs[0] |= 0xE0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 2097152)
				})
				Convey("123456789 (0x075BCD15)", func() {
					binary.BigEndian.PutUint32(b32, 123456789)
					copy(bs, b32[0:4])
					bs[0] |= 0xE0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 123456789)
				})
				Convey("268435455 (0x0FFFFFFF)", func() {
					binary.BigEndian.PutUint32(b32, 268435455)
					copy(bs, b32[0:4])
					bs[0] |= 0xE0
					l, _ := makeApiReader(bs).ReadLen()
					So(l, ShouldEqual, 268435455)
				})
			})
		})
	})
}

func makeApiReader(bs []byte) *ApiReader {
	return NewReader(bytes.NewReader(bs))
}
