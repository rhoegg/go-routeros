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

type Sentence struct {
	Command    string
	Attributes map[string]string
}
type Request struct {
	Sentence
}
type Response struct {
	Done      bool
	Sentences []Sentence
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

func (s *Session) Send(words []string) ([][]string, error) {
	log.Printf(" -->  %v", words)
	var err error
	for _, word := range words {
		_, err = s.writer.WriteWord(word)
	}
	_, err = s.writer.WriteWord("")
	r := [][]string{}
	l, err := s.reader.ReadSentence()
	r = append(r, l)
	for l[0] != "!done" && l[0] != "!trap" && l[0] != "!fatal" {
		log.Printf("  <-- %v", r)
		l, err = s.reader.ReadSentence()
		r = append(r, l)
	}

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

func parseResponse(lines [][]string) Response {
	r := Response{}
	for _, words := range lines {
		r.Sentences = append(r.Sentences,
			Sentence{
				Command:    words[0][1:],
				Attributes: parseAttributes(words)})
	}
	r.Done = (r.Sentences[len(r.Sentences)-1].Command == "done")
	return r
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
