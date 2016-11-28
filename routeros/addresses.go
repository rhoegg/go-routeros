package routeros

import (
    "strconv"
)

type Address struct{
    id string
    Address string
    Interface string
    Dynamic bool
    Invalid bool
    Disabled bool
}

func (s *Session) DescribeAddresses() ([]Address, error) {
    r, err := s.Request(Request{Sentence{
        Command:    "ip/address/print",
        Attributes: map[string]string{}}})
    if err != nil {
        return nil, err
    }
    a := []Address{}
    parsebool := func(s string) bool {
        var b bool
        b, err = strconv.ParseBool(s)
        return b
    }
    for _, s := range r.Sentences {
        if ("re" == s.Command) {
            a = append(a, Address{
                id: s.Attributes[".id"],
                Address: s.Attributes["address"],
                Interface: s.Attributes["interface"],
                Dynamic: parsebool(s.Attributes["dynamic"]),
                Invalid: parsebool(s.Attributes["invalid"]),
                Disabled: parsebool(s.Attributes["disabled"])})
        }
    }
    return a, err
}

func (s *Session) AddAddress(a Address) error {
    _, err := s.Request(Request{Sentence{
        Command: "ip/address/add",
        Attributes: map[string]string{
            "address": a.Address,
            "interface": a.Interface}}})
    return err
}

func (s *Session) RemoveAddress(a Address) error {
    pos := -1
    addresses, err := s.DescribeAddresses()
    if err != nil {
        return err
    }
    for i, address := range addresses {
        if (address.id == a.id || address.Address == a.Address) {
            pos = i
            break
        }
    }
    if pos == -1 {
        // it's already gone
        return nil
    }
    _, err = s.Request(Request{Sentence{
        Command: "ip/address/remove",
        Attributes: map[string]string{
            "numbers": strconv.Itoa(pos)},
        Query: map[string]string{}}})
    return err
}