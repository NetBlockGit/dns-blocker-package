package config

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

var blockerConfig *BlockerConfig

func (bc *BlockerConfig) StartDnsServer() error {
	blockerConfig = bc
	dns.HandleFunc(".", dnsHandler)

	go func() {
		server := &dns.Server{Addr: bc.Addr, Net: "udp"}
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		server := &dns.Server{Addr: bc.Addr, Net: "tcp"}
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	domainName := r.Question[0].Name

	domainNameWithoutDot := strings.TrimSuffix(domainName, ".")
	if blockerConfig.Enabled {
		for _, host := range blockerConfig.BlockList {
			if strings.Contains(domainNameWithoutDot, host) {
				writeNullMsg(&w, r)
				blockerConfig.sendQueryEvent(domainNameWithoutDot, true)
				log.Printf("Blocked %v", domainNameWithoutDot)
				return
			}
		}
	}

	blockerConfig.sendQueryEvent(domainNameWithoutDot, false)
	log.Printf("Passed with domain name %v", domainNameWithoutDot)
	upstreamDNSAddr := blockerConfig.UpstreamDns
	dnsRes, err := dns.Exchange(r, upstreamDNSAddr)
	if err != nil {
		log.Printf("Failed to send message to dns server at %v, reported error - %v", upstreamDNSAddr, err)
		writeNullMsg(&w, r)

		return
	}
	w.WriteMsg(dnsRes)
}

func writeNullMsg(w *dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	m.RecursionAvailable = true
	m.Rcode = dns.RcodeNameError
	(*w).WriteMsg(m)
}
