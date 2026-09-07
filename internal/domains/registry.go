package domains

import (
	"fmt"
	"sync"
)

var (
	mu       sync.Mutex
	registry = make(map[string]Domain)
)

func Register(name string, d Domain) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("domain %q already registered", name))
	}
	registry[name] = d
}

func AllDomains() []Domain {
	mu.Lock()
	defer mu.Unlock()
	domains := make([]Domain, 0, len(registry))
	for _, d := range registry {
		domains = append(domains, d)
	}
	return domains
}

func GetDomain(name string) (Domain, bool) {
	mu.Lock()
	defer mu.Unlock()
	d, ok := registry[name]
	return d, ok
}

func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]Domain)
}
