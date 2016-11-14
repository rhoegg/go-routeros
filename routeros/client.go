package routeros

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedWordLength = errors.New("routeros: word length not supported")
)

type ApiClient struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Session struct {
	Client *ApiClient
	conn   net.Conn
	writer *ApiWriter
	reader *ApiReader
}

type Request struct {
	Command    string
	Attributes map[string]string
}
type Response struct {
	Attributes map[string]string
	words      []string
}

func NewClient(host string, port int, user, password string) *ApiClient {
	return &ApiClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password}
}

func (c *ApiClient) Connect() (*Session, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
	if err != nil {
		return nil, err
	}
	s := &Session{
		Client: c,
		conn:   conn,
		writer: NewWriter(conn),
		reader: NewReader(conn)}
	err = s.Login()
	return s, err
}

func (s *Session) Request(r Request) (Response, error) {
	raw, err := s.Send(request(r))
	return parseResponse(raw), err
}

func (s *Session) Send(words []string) ([]string, error) {
	log.Printf(" -->  %v", words)
	var err error
	for _, word := range words {
		_, err = s.writer.WriteWord(word)
	}
	_, err = s.writer.WriteWord("")
	r, err := s.reader.ReadSentence()
	log.Printf("  <-- %v", r)
	return r, err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

func request(r Request) []string {
	words := []string{fmt.Sprintf("/%s", r.Command)}
	for k, v := range r.Attributes {
		words = append(words, fmt.Sprintf("=%s=%s", k, v))
	}
	return words
}

func parseResponse(words []string) Response {
	return Response{
		Attributes: parseAttributes(words),
		words:      words}
}

func parseAttributes(words []string) map[string]string {
	a := map[string]string{}
	for _, w := range words {
		if w[0] == '=' {
			parts := strings.SplitN(w[1:], "=", 2)
			a[parts[0]] = parts[1]
		}
	}
	return a
}
