package routeros

import (
    "errors"
    "fmt"
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