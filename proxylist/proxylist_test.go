package proxylist

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var feed = `
[
  {"host":"177.207.216.238","ip":"177.207.216.238","port":"3128","lastseen":60,"delay":4020,"cid":"48220","country_code":"BR","country_name":"Brazil","city":"Brusque","checks_up":"1","checks_down":"0","anon":"1","http":"1","ssl":"0","socks4":"0","socks5":"0"},
  {"host":"173.214.148.113","ip":"173.214.148.113","port":"10200","lastseen":60,"delay":1580,"cid":"13700","country_code":"US","country_name":"United States","city":"Cairo","checks_up":"7","checks_down":"29","anon":"4","http":"0","ssl":"0","socks4":"1","socks5":"0"}
]
`

var notFound = `NOTFOUND`

var tooFast = `TOOFAST`

func mockServer(body string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=windows-1251")
		fmt.Fprintln(w, body)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestValidLoad(t *testing.T) {
	server := mockServer(feed)
	defer server.Close()
	ps, err := Load(server.URL)

	assert.Equal(t, 2, len(ps))
	assert.NoError(t, err)
	p := ps[0]
	assert.Equal(t, 3128, p.Port)
	assert.True(t, bool(p.HTTP))
	assert.False(t, bool(p.SSL))
	assert.Equal(t, "http://177.207.216.238:3128", p.ToURL())

	p = ps[1]
	assert.Equal(t, 10200, p.Port)
	assert.False(t, bool(p.HTTP))
	assert.True(t, bool(p.Socks4))
	assert.Equal(t, "socks4://173.214.148.113:10200", p.ToURL())
}

func TestWrongCodeLoad(t *testing.T) {
	server := mockServer(notFound)
	defer server.Close()
	ps, err := Load(server.URL)
	assert.Nil(t, ps)
	assert.Error(t, err)
}

func TestTooFastLoad(t *testing.T) {
	server := mockServer(tooFast)
	defer server.Close()
	ps, err := Load(server.URL)
	assert.Nil(t, ps)
	assert.Error(t, err)
}
