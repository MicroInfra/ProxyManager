package models

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"syscall"
)

type Proxy struct {
	ServiceName string `json:"serviceName"` // Unique name of the service
	ServiceUrl  string `json:"serviceUrl"`  // The URL of the service in Docker network. Must end with /
	ListenPort  int    `json:"listenPort"`  // The port to listen on
	ProxyType   string `json:"proxyType"`   // The type of proxy to use. [http, tcp]
	FilterFile  string `json:"filterFile"`  // The file containing the filter rules. Path of file in future
	Pid         int
}

type ProxyRulesPlainText struct {
	ServiceName string `json:"serviceName"` // Unique name of the service
	ServiceUrl  string `json:"serviceUrl"`  // The URL of the service in Docker network. Must end with /
	ListenPort  int    `json:"listenPort"`  // The port to listen on
	ProxyType   string `json:"proxyType"`   // The type of proxy to use. [http, tcp]
	Rules       string `json:"rules"`       // The code of rules
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
		cmd = exec.Command("mitmdump", "--set", "block_global=false", "--mode", fmt.Sprintf("reverse:%s", proxy.ServiceUrl), "-s", proxy.FilterFile, "--listen-port", fmt.Sprint(proxy.ListenPort), "&")
	} else {
		cmd = exec.Command("mitmdump", "--set", "block_global=false", "--mode", fmt.Sprintf("reverse:%s", proxy.ServiceUrl), "--listen-port", fmt.Sprint(proxy.ListenPort), "&")
	}

	if err := cmd.Start(); err != nil {
		log.Println("Error starting command")
		return fmt.Errorf("Error starting command: %s\n", err)
	}
	proxy.Pid = cmd.Process.Pid
	log.Println("Proxy PID is ", proxy.Pid)
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
	pid := p.Proxies[name].Pid
	if pid == 0 {
		return
	}
	defer p.Unlock()
	p.Lock()
	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		fmt.Printf("Error killing command: %s\n", err)
	}

	delete(p.Proxies, name)
}
