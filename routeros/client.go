package routeros

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
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
	Client     *ApiClient
	conn       net.Conn
	writer     *ApiWriter
	reader     *ApiReader
	cmdCounter int64
	replies    map[string]chan Sentence
}

type item interface {
	toAttributes() map[string]string
}
type mapItem struct {
	Map map[string]string
}

func (m mapItem) toAttributes() map[string]string {
	return m.Map
}
func itemFromMap(m map[string]string) item {
	return mapItem{m}
}
func emptyItem() item {
	return itemFromMap(map[string]string{})
}

type Sentence struct {
	Command       string
	Attributes    map[string]string
	Query         map[string]string
	ApiAttributes map[string]string
	words         []string
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
		Client:  c,
		conn:    conn,
		writer:  NewWriter(conn),
		reader:  NewReader(conn),
		replies: map[string]chan Sentence{}}
	err = s.Login()
	if err == nil {
		go s.Receive()
	}
	return s, err
}

func (s *Session) Request(c string, i item) (*Response, error) {
	tag := strconv.FormatInt(atomic.AddInt64(&s.cmdCounter, 1), 10)

	r := request(Request{Sentence{
		Command:       c,
		Attributes:    i.toAttributes(),
		ApiAttributes: map[string]string{"tag": tag}}})
	// buffer up to 32 reply sentences for a single tag before blocking
	s.replies[tag] = make(chan Sentence, 32)
	if err := s.Send(r); err == nil {
		return parseResponse(s.replies[tag]), err
	} else {
		return nil, err
	}
}

func (s *Session) Send(words []string) error {
	log.Printf(" -->  %v", words)
	var err error
	for _, word := range words {
		_, err = s.writer.WriteWord(word)
	}
	_, err = s.writer.WriteWord("")
	return err
}

func (s *Session) Receive() {
	for {
		words, err := s.reader.ReadSentence()
		log.Printf("Received a sentence %v", words)
		if err != nil {
			log.Printf("Bad news! %v", err)
			return
		}
		s.receiveSentence(parseSentence(words))
	}
}

func (s *Session) receiveSentence(sentence Sentence) {
	tag, tagIncluded := sentence.ApiAttributes["tag"]
	if tagIncluded {
		replyChannel, tagExpected := s.replies[tag]
		if tagExpected {
			log.Printf("  <-- ( %s) %v", tag, sentence.words)
			replyChannel <- sentence
			if sentence.Command == "done" {
				close(replyChannel)
			}
		} else {
			log.Printf("  <-- (?%s) %v", tag, sentence.words)
		}
	} else {
		log.Printf("  <-- (xx) %v", sentence.words)
	}
}

func (s *Session) Close() error {
	return s.conn.Close()
}

func request(r Request) []string {
	words := []string{fmt.Sprintf("/%s", r.Command)}
	for k, v := range r.ApiAttributes {
		words = append(words, fmt.Sprintf(".%s=%s", k, v))
	}
	for k, v := range r.Attributes {
		words = append(words, fmt.Sprintf("=%s=%s", k, v))
	}
	for k, v := range r.Query {
		words = append(words, fmt.Sprintf("?%s=%s", k, v))
	}
	return words
}

func parseResponse(lines <-chan Sentence) *Response {
	r := Response{}
	for s := range lines {
		r.Sentences = append(r.Sentences, s)
	}
	r.Done = (r.Sentences[len(r.Sentences)-1].Command == "done")
	return &r
}

func parseSentence(words []string) Sentence {
	return Sentence{
		Command:       words[0][1:],
		ApiAttributes: parseApiAttributes(words),
		Attributes:    parseAttributes(words),
		Query:         parseQuery(words),
		words:         words}
}

func parseApiAttributes(words []string) map[string]string {
	a := map[string]string{}
	for _, w := range words {
		if w[0] == '.' {
			parts := strings.SplitN(w[1:], "=", 2)
			a[parts[0]] = parts[1]
		}
	}
	return a
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

func parseQuery(words []string) map[string]string {
	a := map[string]string{}
	for _, w := range words {
		if w[0] == '?' {
			parts := strings.SplitN(w[1:], "=", 2)
			a[parts[0]] = parts[1]
		}
	}
	return a
}
