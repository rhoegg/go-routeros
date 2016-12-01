package routeros

import (
	"github.com/fatih/structs"
)

type IpFirewallFilterRule struct {
	Chain          string `structs:"chain"`
	Action         string `structs:"action"`
	Protocol       string `structs:"protocol,omitempty"`
	RejectWith     string `structs:"reject-with,omitempty"`
	SrcAddressList string `structs:"src-address-list,src-address-list"`
}

func (r IpFirewallFilterRule) toAttributes() map[string]string {
	a := map[string]string{}
	for k, v := range structs.Map(r) {
		a[k] = v.(string)
	}
	return a
}

func (s *Session) DescribeIpFirewallFilter() ([]IpFirewallFilterRule, error) {
	r, err := s.Request("ip/firewall/filter/print", emptyItem())
	if err != nil {
		return nil, err
	}
	rules := []IpFirewallFilterRule{}
	for _, s := range r.Sentences {
		if "re" == s.Command {
			rules = append(rules, IpFirewallFilterRule{
				Chain:          s.Attributes["chain"],
				Action:         s.Attributes["action"],
				Protocol:       s.Attributes["protocol"],
				RejectWith:     s.Attributes["reject-with"],
				SrcAddressList: s.Attributes["src-address-list"]})
		}
	}
	return rules, err
}

func (s *Session) AddIpFirewallFilterRule(r IpFirewallFilterRule) error {
	_, err := s.Request("ip/firewall/filter/add", r)
	return err
}
