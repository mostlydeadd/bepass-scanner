package scanner

import (
	"github.com/kaveh-ahangar/cfscanner/internal/config"
	"sync"
)

type Scanner struct {
	activeIPs []string
	sync.Mutex
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Run() {
	// Use channels to collect results from goroutines.

	convertCIDRtoIPList()

	var wg sync.WaitGroup
	activeIPChan := make(chan string, config.G.Threads*2)
	defer close(activeIPChan)

	for i := 0; i < config.G.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range activeIPChan {
				s.Lock()
				s.activeIPs = append(s.activeIPs, ip)
				s.Unlock()
			}
		}()
	}

	if config.G.PingMode {
		pingTechnique(activeIPChan)
	}

	if config.G.PortscanMode {
		portScanTechnique(activeIPChan)
	}

	wg.Wait()
}
