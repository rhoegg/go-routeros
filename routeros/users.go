package routeros

import (
	"strconv"
)

type User struct {
	id       string
	Name     string
	Group    string
	Password string
	Address  string
	Comment  string
	Disabled bool
}

func (s *Session) DescribeUsers() ([]User, error) {
	r, err := s.Request(Request{Sentence{
		Command:    "user/print",
		Attributes: map[string]string{}}})
	if err != nil {
		return nil, err
	}
	u := []User{}
	for _, s := range r.Sentences {
		if "re" == s.Command {
			u = append(u, User{
				id:       s.Attributes[".id"],
				Name:     s.Attributes["name"],
				Group:    s.Attributes["group"],
				Password: s.Attributes["password"],
				Address:  s.Attributes["address"],
				Comment:  s.Attributes["comment"],
				Disabled: parseBool(s.Attributes["disabled"])})
		}
	}
	return u, err
}

func (s *Session) AddUser(u User) error {
	_, err := s.Request(Request{Sentence{
		Command: "user/add",
		Attributes: map[string]string{
			"name":     u.Name,
			"password": u.Password,
			"group":    u.Group,
			"address":  u.Address,
			"comment":  u.Comment}}})
	return err
}

func (s *Session) RemoveUser(u User) error {
	return s.withUserIndex(u, func(pos int) error {
		_, err := s.Request(Request{Sentence{
			Command: "user/remove",
			Attributes: map[string]string{
				"numbers": strconv.Itoa(pos)}}})
		return err
	})
}

func (s *Session) withUserIndex(u User, action func(int) error) error {
	users, err := s.DescribeUsers()
	if err != nil { return err }
	for i, user := range users {
		if user.id == u.id || user.Name == u.Name {
			return action(i)
		}
	}
	return nil
}
