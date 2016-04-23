package proxylist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type boolFlag bool

func (bf *boolFlag) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	if s == "1" {
		*bf = true
	} else {
		*bf = false
	}
	return nil
}

// Proxy type
type Proxy struct {
	Host        string   `json:"host"`
	IP          net.IP   `json:"ip"`
	Port        int      `json:"port,string"`
	LastSeen    int      `json:"lastseen"`
	Delay       int      `json:"delay"`
	CID         int      `json:"cid,string"`
	CountryCode string   `json:"country_code"`
	CountryName string   `json:"country_name"`
	City        string   `json:"city"`
	ChecksUp    int      `json:"checks_up,string"`
	ChecksDown  int      `json:"checks_down,string"`
	Anon        int      `json:"anon,string"`
	HTTP        boolFlag `json:"http"`
	SSL         boolFlag `json:"ssl"`
	Socks4      boolFlag `json:"socks4"`
	Socks5      boolFlag `json:"socks5"`
}

// ToURL makes proxy URL
func (proxy *Proxy) ToURL() string {
	protocol := "http"

	if proxy.SSL {
		protocol = "https"
	}

	if proxy.Socks4 {
		protocol = "socks4"
	}

	if proxy.Socks5 {
		protocol = "socks5"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, proxy.IP, proxy.Port)
}

// Load proxy list
func Load(uri string) ([]Proxy, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(buf) == "NOTFOUND" {
		return nil, errors.New("Wrong code provided")
	}
	if string(buf) == "TOOFAST" {
		return nil, errors.New("Frequency restriction met")
	}
	var p []Proxy
	err = json.Unmarshal(buf, &p)
	if err != nil {
		return nil, err
	}
	return p, err
}
