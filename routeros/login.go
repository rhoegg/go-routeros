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
	r, err := s.Request("login", emptyItem())
	if !r.Done {
		return ErrUnsuccessfulLoginResult
	}
	challenge := r.Sentences[0].Attributes["ret"]
	response, err := encodePassword(s.Client.Password, challenge)
	r, err = s.Request(
		"login",
		itemFromMap(map[string]string{
			"name":     s.Client.User,
			"response": response}))
	if !r.Done {
		return ErrUnsuccessfulLoginResult
	}
	return err
}

func (s *Session) Quit() error {
	_, err := s.Request("quit", emptyItem())
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
