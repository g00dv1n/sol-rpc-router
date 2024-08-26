package router

import (
	"sync"
	"testing"
)

const (
	totalReq = 100000
	concrReq = 25
)

func TestRoundRobinBalancer(t *testing.T) {
	servers := []ServerEndpoint{
		MustConvertToEndpoint("http://serverA.com", 0),
		MustConvertToEndpoint("http://serverB.com", 0),
		MustConvertToEndpoint("http://serverC.com", 0),
	}

	balancer, _ := NewRoundRobinBalancer(servers)

	serverCounts := BalancerTest(balancer, totalReq, concrReq)

	t.Logf("Total reqs: %d | Concurrent: %d", totalReq, concrReq)
	t.Logf("Reqs by server: %v", serverCounts)
}

func TestWeightedRoundRobinBalancer(t *testing.T) {
	servers := []ServerEndpoint{
		MustConvertToEndpoint("http://serverA.com", 1),
		MustConvertToEndpoint("http://serverB.com", 3),
		MustConvertToEndpoint("http://serverC.com", 1),
	}

	balancer, _ := NewWeightedRoundRobinBalancer(servers)

	serverCounts := BalancerTest(balancer, totalReq, concrReq)

	t.Logf("Total reqs: %d | Concurrent: %d", totalReq, concrReq)
	t.Logf("Reqs by server: %v", serverCounts)
}

type ServerCounts map[string]int

func BalancerTest(balancer Balancer, totalReq int, concrReq int) map[string]int {
	serverCounts := make(ServerCounts)
	countsLock := sync.Mutex{}

	for i := 0; i < totalReq; {
		wg := sync.WaitGroup{}
		for j := 0; j < concrReq && i < totalReq; j++ {
			wg.Add(1)
			go func() {
				countsLock.Lock()
				defer wg.Done()
				defer countsLock.Unlock()

				s := balancer.NextServer()
				serverCounts[s.Host]++

			}()
			i++
		}
		wg.Wait()
	}

	return serverCounts
}
