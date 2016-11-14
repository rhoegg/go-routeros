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
	r, err := s.Request(Request{"login", map[string]string{}})
	if r.words[0] != "!done" {
		return ErrUnsuccessfulLoginResult
	}
	challenge := r.Attributes["ret"]
	response, err := encodePassword(s.Client.Password, challenge)
	r, err = s.Request(Request{
		"login",
		map[string]string{
			"name":     s.Client.User,
			"response": response}})
	if r.words[0] != "!done" {
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
