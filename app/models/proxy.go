package models

import (
  "fmt"
	"sync"
  "os/exec"
  "syscall"
)


type Proxy struct {
	ServiceName string    `json:"serviceName"` // Unique name of the service
	ServiceUrl  string    `json:"serviceUrl"`  // The URL of the service in Docker network. Must end with /
	ListenPort  int       `json:"listenPort"`  // The port to listen on
	ProxyType   string    `json:"proxyType"`   // The type of proxy to use. [http, tcp]
	FilterFile  string    `json:"filterFile"`  // The file containing the filter rules
  Pid         int
}

type Proxies struct {
	Proxies map[string]Proxy `json:"proxies"`
	sync.Mutex
}

func NewAllProxies() Proxies {
	return Proxies{Proxies: make(map[string]Proxy)}
}

func (p *Proxies) Set(name string, proxy Proxy) error {
	defer p.Unlock()
	p.Lock()
	p.Proxies[name] = proxy
  var cmd *exec.Cmd 
  if proxy.FilterFile != "" {
    cmd = exec.Command("mitmdump", "--set", "block_global=false", "--mode", fmt.Sprintf("reverse:%s", proxy.ServiceUrl), "-s", proxy.FilterFile, "&")

  } else {
    cmd = exec.Command("mitmdump", "--set", "block_global=false", "--mode", fmt.Sprintf("reverse:%s", proxy.ServiceUrl), "&")

  }

  if err := cmd.Start(); err != nil {
    return fmt.Errorf("Error starting command: %s\n", err)
  }
  proxy.Pid = cmd.Process.Pid
  p.Proxies[name] = proxy

  return nil
}

func (p *Proxies) Get(name string) Proxy {
	defer p.Unlock()
	p.Lock()
	val := p.Proxies[name]
	return val
}


func (p *Proxies) Delete(name string) {
	defer p.Unlock()
	p.Lock()
  if err := syscall.Kill(p.Proxies[name].Pid, syscall.SIGKILL); err != nil {
      fmt.Printf("Error killing command: %s\n", err)
  }

	delete(p.Proxies, name)
}
