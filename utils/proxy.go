package utils

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
)

type Proxy struct {
	Raw      string
	Scheme   string
	Host     string
	Port     string
	Username string
	Password string
}

func (p *Proxy) URL() string {
	if p.Username != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s", p.Scheme, p.Username, p.Password, p.Host, p.Port)
	}
	return fmt.Sprintf("%s://%s:%s", p.Scheme, p.Host, p.Port)
}

type ProxyPool struct {
	proxies []Proxy
	idx     atomic.Int64
}

func LoadProxies(filePath string) (*ProxyPool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open proxy file: %w", err)
	}
	defer f.Close()

	pool := &ProxyPool{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if p, err := parseProxy(line); err == nil {
			pool.proxies = append(pool.proxies, p)
		}
	}
	if len(pool.proxies) == 0 {
		return nil, fmt.Errorf("no valid proxies found in %s", filePath)
	}
	return pool, scanner.Err()
}

func (p *ProxyPool) Next() string {
	if len(p.proxies) == 0 {
		return ""
	}
	idx := int(p.idx.Add(1)-1) % len(p.proxies)
	return p.proxies[idx].URL()
}

func (p *ProxyPool) Len() int {
	return len(p.proxies)
}

func parseProxy(raw string) (Proxy, error) {

	target := raw
	if !strings.Contains(raw, "://") {
		target = "http://" + raw
	}
	u, err := url.Parse(target)
	if err != nil {
		return Proxy{}, err
	}
	p := Proxy{
		Raw:    raw,
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Port:   u.Port(),
	}
	if u.User != nil {
		p.Username = u.User.Username()
		p.Password, _ = u.User.Password()
	}
	return p, nil
}
