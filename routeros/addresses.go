package routeros

import (
	"github.com/fatih/structs"
	"strconv"
)

type Address struct {
	id        string
	Address   string `structs:"address"`
	Interface string `structs:"interface"`
	Dynamic   bool   `structs:"dynamic,omitempty"`
	Invalid   bool   `structs:"invalid,omitempty"`
	Disabled  bool   `structs:"disabled,omitempty"`
}

func (a Address) toAttributes() map[string]string {
	attrs := map[string]string{}
	for k, v := range structs.Map(a) {
		attrs[k] = v.(string)
	}
	return attrs
}

func (s *Session) DescribeAddresses() ([]Address, error) {
	r, err := s.Request("ip/address/print", emptyItem())
	if err != nil {
		return nil, err
	}
	a := []Address{}
	for _, s := range r.Sentences {
		if "re" == s.Command {
			a = append(a, Address{
				id:        s.Attributes[".id"],
				Address:   s.Attributes["address"],
				Interface: s.Attributes["interface"],
				Dynamic:   parseBool(s.Attributes["dynamic"]),
				Invalid:   parseBool(s.Attributes["invalid"]),
				Disabled:  parseBool(s.Attributes["disabled"])})
		}
	}
	return a, err
}

func (s *Session) AddAddress(a Address) error {
	_, err := s.Request("ip/address/add", a)
	return err
}

func (s *Session) RemoveAddress(a Address) error {
	return s.withAddressIndex(a, func(pos int) error {
		_, err := s.Request("ip/address/remove",
			itemFromMap(map[string]string{
				"numbers": strconv.Itoa(pos)}))
		return err
	})
}

func (s *Session) withAddressIndex(a Address, action func(int) error) error {
	addresses, err := s.DescribeAddresses()
	if err != nil {
		return err
	}
	for i, address := range addresses {
		if address.id == a.id || address.Address == a.Address {
			return action(i)
		}
	}
	return nil
}
