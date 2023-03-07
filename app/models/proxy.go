package models

import (
	"net/http"
	"sync"
)


type Proxy struct {
	ServiceName string    `json:"serviceName"` // Unique name of the service
	ServiceUrl  string    `json:"serviceUrl"`  // The URL of the service in Docker network. Must end with /
	ListenPort  int       `json:"listenPort"`  // The port to listen on
	ProxyType   string    `json:"proxyType"`   // The type of proxy to use. [http, tcp]
	FilterFile  http.File `json:"filterFile"`  // The file containing the filter rules
}

type AllProxies struct {
	Proxies map[string]Proxy `json:"proxies"`
	sync.Mutex
}

func NewAllProxies() AllProxies {
	return AllProxies{Proxies: make(map[string]Proxy)}
}

func (p *AllProxies) Set(name string, proxy Proxy) {
	defer p.Unlock()
	p.Lock()
  print("Vsem privet kto smotrit moi kanal")
	p.Proxies[name] = proxy
}

func (p *AllProxies) Get(name string) Proxy {
	p.Lock()
	val := p.Proxies[name]
	defer p.Unlock()
	return val
}
