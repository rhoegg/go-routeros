package routeros

import (
	"github.com/fatih/structs"
	"strconv"
)

type User struct {
	id       string
	Name     string `structs:"name"`
	Group    string `structs:"group"`
	Password string `structs:"password,omitempty"`
	Address  string `structs:"address,omitempty"`
	Comment  string `structs:"comment,omitempty"`
	Disabled bool   `structs:"disabled,omitempty"`
}

func (u User) toAttributes() map[string]string {
	attrs := map[string]string{}
	for k, v := range structs.Map(u) {
		attrs[k] = v.(string)
	}
	return attrs
}

func (s *Session) DescribeUsers() ([]User, error) {
	r, err := s.Request("user/print", emptyItem())
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
	_, err := s.Request("user/add", u)
	return err
}

func (s *Session) RemoveUser(u User) error {
	return s.withUserIndex(u, func(pos int) error {
		_, err := s.Request("user/remove",
			itemFromMap(map[string]string{
				"numbers": strconv.Itoa(pos)}))
		return err
	})
}

func (s *Session) withUserIndex(u User, action func(int) error) error {
	users, err := s.DescribeUsers()
	if err != nil {
		return err
	}
	for i, user := range users {
		if user.id == u.id || user.Name == u.Name {
			return action(i)
		}
	}
	return nil
}
