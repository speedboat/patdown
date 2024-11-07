package common

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

type Query struct {
	Nameserver string
	Vendor     string
	DomainPair Pair
}

type Nameserver struct {
	Nameserver string
	NonRA      bool
	Recursive  bool
}

func message(domain string, reqtype uint16, ra bool) *dns.Msg {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = ra
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  reqtype,
		Qclass: dns.ClassINET,
	}
	return msg
}

func ParseNS(nameservers []string) []Nameserver {
	var valid []Nameserver
	msg := message("cloudflare.com", dns.TypeA, false)
	for _, ns := range nameservers {
		nonra, ra := false, false
		in, err := dns.Exchange(msg, ns+":53")
		if err != nil {
			Error(fmt.Sprintf("nameserver %s%s%s is not responding to the trial query", ColorGray, ns[0:len(ns)-1], ColorReset))
			continue
		}
		if in.Rcode == dns.RcodeRefused {
			Warn(fmt.Sprintf("nameserver %s%s%s refused the trial non-recursive query", ColorGray, ns[0:len(ns)-1], ColorReset))
		} else {
			Success(fmt.Sprintf("nameserver %s%s%s allows non-recursive queries", ColorGray, ns[0:len(ns)-1], ColorReset))
			nonra = true
		}
		if in.RecursionAvailable {
			Success(fmt.Sprintf("nameserver %s%s%s allows recursion", ColorGray, ns[0:len(ns)-1], ColorReset))
			ra = true
		} else {
			Warn(fmt.Sprintf("nameserver %s%s%s does not allow recursion", ColorGray, ns[0:len(ns)-1], ColorReset))
		}

		valid = append(valid, Nameserver{Nameserver: ns, NonRA: nonra, Recursive: ra})
	}
	return valid
}

func NeutralReq() bool {
	msg := message("supernets.org", dns.TypeA, true)
	in, err := dns.Exchange(msg, "1.1.1.1:53")
	if err != nil {
		return false
	}
	if len(in.Answer) > 0 {
		return true
	}
	return false
}

func PullNS(d string) []string {
	nsmsg := message(d, dns.TypeNS, true)
	in, err := dns.Exchange(nsmsg, "1.1.1.1:53")
	if err != nil {
		Fatal("unable to retrieve nameservers for " + d)
	}

	nameservers := []string{}

	for _, ans := range in.Answer {
		ns, ok := ans.(*dns.NS)
		if ok {
			nameservers = append(nameservers, ns.Ns)
		}
	}

	return nameservers
}

func RunQuery(q <-chan Query, tracker chan<- interface{}, delay int) {
	for qdata := range q {
		if Params.Verbose {
			Info(fmt.Sprintf("querying %s on %s", qdata.DomainPair.Domain, qdata.Nameserver[0:len(qdata.Nameserver)-1]))
		}
		msg := message(qdata.DomainPair.Domain, dns.TypeA, false)
		in, err := dns.Exchange(msg, qdata.Nameserver+":53")
		if err != nil {
			Error(err.Error())
			continue
		}

		if len(in.Answer) > 0 {
			Success(fmt.Sprintf("[%s] associated domain %s%s%s found on %s%s%s",
				qdata.Vendor, ColorRed, qdata.DomainPair.Domain, ColorReset, ColorRed, qdata.Nameserver[0:len(qdata.Nameserver)-1], ColorReset))
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	tracker <- 1337
}

func RunQueryRA(q <-chan Query, tracker chan<- interface{}, delay int) {
	for qdata := range q {
		if Params.Verbose {
			Info(fmt.Sprintf("recursively querying %s on %s", qdata.DomainPair.Domain, qdata.Nameserver[0:len(qdata.Nameserver)-1]))
		}
		for x := 0; x < 2; x++ {
			msg := message(qdata.DomainPair.Domain, dns.TypeA, true)
			in, err := dns.Exchange(msg, qdata.Nameserver+":53")
			if err != nil {
				Error("hiccup on " + qdata.Nameserver[0:len(qdata.Nameserver)-1] + " while querying " + qdata.DomainPair.Domain)
				time.Sleep(2 * time.Second)
				continue
			}

			if len(in.Answer) > 0 {
				if in.Answer[0].Header().Ttl <= qdata.DomainPair.TTL-4 {
					Success(fmt.Sprintf("[%s] associated domain %s%s%s found on %s%s%s with decremented TTL of %s%d%s",
						qdata.Vendor, ColorRed, qdata.DomainPair.Domain, ColorReset, ColorRed, qdata.Nameserver[0:len(qdata.Nameserver)-1], ColorReset, ColorGreen, in.Answer[0].Header().Ttl, ColorReset))
				}
			}
			break
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	tracker <- 1337
}
