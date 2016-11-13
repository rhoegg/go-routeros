package routeros

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net"
    "strconv"
)

var (
    ErrUnsupportedWordLength = errors.New("routeros: word length not supported")
)

type ApiClient struct {
    Host string
    Port int
    User string
    Password string
}

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

// Algorithm for length is at http://wiki.mikrotik.com/wiki/Manual:API#Protocol
func (r *ApiReader) ReadLen() (uint32, error) {
    threshold := [...]byte{0, 0x80, 0xC0, 0xE0, 0xF0}

    first, err := r.rd.ReadByte()
    log.Printf("first byte (%d)", first)
    i := uint32(first)
    for t := 0; t<len(threshold); t++ {
        if first < threshold[t + 1] {
            // 0x00, 0x8000, 0xC00000, 0xE0000000
            mask := uint32(threshold[t]) << uint32(t * 8)
            return i &^ mask, err
        }
        b, _ := r.rd.ReadByte()
        i = i << 8 + uint32(b)
        log.Printf("byte %d: %d", t + 2, b)
    }
    return i, ErrUnsupportedWordLength
}

func (c *ApiClient) Talk(words []string) []string {
    conn, err := net.Dial("tcp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
    if (err != nil) {
        // oops
        log.Printf("Something's borked: %v", err)
    }
    defer conn.Close()
    log.Printf(" -->  %v", words)
    for _, word := range words {
        fmt.Fprintf(conn, word)
    }
    response := words
    log.Printf("  <-- %v", response)
    return words
}