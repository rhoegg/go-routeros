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
    data, err := ioutil.ReadAll(io.LimitReader(r.rd, l))
    return string(data), err
}

func (r *ApiReader) ReadLen() (int64, error) {
    b, err := r.rd.ReadByte()
    if b > byte(0x7F) {
        return int64(b), ErrUnsupportedWordLength
    }
    return int64(b), err
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