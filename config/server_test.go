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
		BlockList:   []string{},
		Addr:        addr,
		Enabled:     true,
	}
	go bc.StartDnsServer()
	time.Sleep(2 * time.Second)

	t.Run("Should block if it is present in blocklist", func(t *testing.T) {
		bc.AddHostToBlockList("ommore.me")
		r := tryAQuery(t, "ommore.me", addr)
		assert.Len(t, r.Answer, 0)
	})

	t.Run("Should respond with answer if it is not present in blocklist", func(t *testing.T) {
		bc.UpdateBlockList([]string{})
		r := tryAQuery(t, "ommore.me", addr)
		assert.Len(t, r.Answer, 2)
	})

	t.Run("Should not block if blocker is disabled", func(t *testing.T) {
		bc.Enabled = false
		bc.AddHostToBlockList("ommore.me")
		r := tryAQuery(t, "ommore.me", addr)
		assert.Len(t, r.Answer, 2)
	})

	t.Run("should receive query if channel is set", func(t *testing.T) {
		ch := make(chan QueryEvent)
		bc.QueryChannel = ch
		go func() {
			for res := range ch {
				assert.Equal(t, QueryEvent{Hostname: "ommore.me", Blocked: false}, res)
				break
			}
		}()
		r := tryAQuery(t, "ommore.me", addr)
		assert.Len(t, r.Answer, 2)

	})

}

func tryAQuery(t *testing.T, host string, addr string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(host+".", dns.TypeA)
	r, e := dns.Exchange(msg, addr)
	if e != nil {
		t.Fatal(e)
	}
	return r
}
