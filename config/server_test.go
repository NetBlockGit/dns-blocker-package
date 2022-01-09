package config

import (
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func Test_DNSBlocker(t *testing.T) {
	addr := "127.0.0.1:8000"
	bc := BlockerConfig{
		UpstreamDns: "1.1.1.1:53",
		BlockList:   []string{"ommore.me"},
		Addr:        addr,
	}
	go bc.StartDnsServer()
	time.Sleep(2 * time.Second)

	t.Run("Should block if it is present in blocklist", func(t *testing.T) {
		msg := new(dns.Msg)
		msg.SetQuestion("ommore.me.", dns.TypeA)
		r, e := dns.Exchange(msg, addr)
		if e != nil {
			t.Fatal(e)
		}
		assert.Len(t, r.Answer, 0)
	})

	t.Run("Should respond with answer if it is not present in blocklist", func(t *testing.T) {
		bc.UpdateBlockList([]string{})
		msg := new(dns.Msg)
		msg.SetQuestion("ommore.me.", dns.TypeA)
		r, e := dns.Exchange(msg, addr)
		if e != nil {
			t.Fatal(e)
		}
		assert.Len(t, r.Answer, 2)
	})

}
