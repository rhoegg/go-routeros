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
    first, err := r.rd.ReadByte()
    log.Printf("we got to level 1 (%d)", first)
    i := uint32(first)
    if first < byte(0x80) {
        return i, err
    }
    b, err := r.rd.ReadByte()
    i = i << 8 + uint32(b)
    log.Printf("we got to level 2 (%d)", b)
    if first < byte(0xC0) {
        return i &^ 0x8000, err
    }
    b, err = r.rd.ReadByte()
    i = i << 8 + uint32(b)
    log.Printf("we got to level 3 (%d)", b)
    if first < byte(0xE0) {
        return i &^ 0xC00000, err
    }
    b, err = r.rd.ReadByte()
    i = i << 8 + uint32(b)
    log.Printf("we got to level 4 (%d)", b)
    if first < byte(0xF0) {
        return i &^ 0xE0000000, err
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