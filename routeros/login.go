package routeros

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

var (
	ErrUnsuccessfulLoginResult = errors.New("routeros: Unrecognized result from /login")
)

func (s *Session) Login() error {
	if err := s.Send(request(Request{Sentence{Command: "login"}})); err != nil {
		return err
	}
	words, err := s.reader.ReadSentence()
	sentence := parseSentence(words)
	challenge := sentence.Attributes["ret"]
	response, err := encodePassword(s.Client.Password, challenge)
	if err != nil {
		return err
	}
	err = s.Send(request(Request{Sentence{
		Command: "login",
		Attributes: map[string]string{
			"name":     s.Client.User,
			"response": response}}}))
	if err != nil {
		return err
	}
	words, err = s.reader.ReadSentence()
	sentence = parseSentence(words)
	if sentence.Command != "done" {
		return ErrUnsuccessfulLoginResult
	}
	return err
}

func encodePassword(p string, c string) (string, error) {
	hash, err := hex.DecodeString(c)
	response := []byte{0}
	response = append(response, []byte(p)...)
	response = append(response, hash...)
	r := md5.Sum(response)
	return "00" + hex.EncodeToString(r[:]), err
}
